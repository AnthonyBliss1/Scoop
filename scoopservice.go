package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yosssi/gohtml"
	"moul.io/http2curl"
)

type Method string

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Patch  Method = "PATCH"
	Delete Method = "DELETE"
	Empty  Method = ""
)

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Method  Method `json:"method"`
	URL     string `json:"url"`
	Headers []KV   `json:"headers"`
	QParams []KV   `json:"query_params"`
	Body    string `json:"body"`
}

type Response struct {
	Status      string `json:"status"`
	StatusCode  int    `json:"status_code"`
	Headers     []KV   `json:"headers"`
	Body        string `json:"body"`
	Duration    int64  `json:"duration"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

type Scoop struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Collection struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Scoops []Scoop `json:"scoops"`
}

type DNSOverride struct {
	Variable string `json:"variable"`
	IPV4     string `json:"ipv4"`
}

type ScoopService struct {
	context context.Context
}

// Initializes the Scoop model by attaching method, url, headers, query params, and body

func (b *ScoopService) ModelIntializer(method Method, reqURL string, headers []KV, qParams []KV, body string) (Request, error) {
	var r Request

	r.Method = method
	r.URL = reqURL
	r.Headers = headers
	r.QParams = qParams
	r.Body = body

	return r, nil
}

// TODO:
// These two funcs - AddQueryParams() and CreateRequest() dont need to be implemented on ScoopService
// since they're purely backend functions
// need to clean up the rest of this file and separate it into 'utils'

func AddQueryParams(s *Scoop) error {
	u, err := url.Parse(s.Request.URL)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return err
	}
	query := url.Values{}

	for _, param := range s.Request.QParams {
		if param.Key == "" || param.Value == "" {
			continue
		}

		query.Add(param.Key, param.Value)
	}

	u.RawQuery = query.Encode()

	s.Request.URL = u.String()

	return nil
}

func CreateRequest(s Scoop, url string) (r *http.Request, err error) {
	if s.Request.Body != "" {
		reqBody := strings.NewReader(s.Request.Body)

		r, err = http.NewRequest(string(s.Request.Method), url, reqBody)
		if err != nil {
			return nil, err
		}

	} else {
		r, err = http.NewRequest(string(s.Request.Method), url, nil)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func RemoveCollectionDir() error {
	base, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	return os.RemoveAll(scoopDir)
}

func (b *ScoopService) SubmitRequest(s Scoop) {
	go func() {
		var r Response

		client := http.Client{}

		// add query params to url
		AddQueryParams(&s)

		// check for use of DNS Overrides
		realURL, err := b.CheckDNSOverride(s)
		if err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return
		}

		// incase no matches
		if realURL == "" {
			realURL = s.Request.URL
		}

		// if there is a body string then need to add it to the request
		req, err := CreateRequest(s, realURL)
		if err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return
		}

		// add headers to request
		for _, h := range s.Request.Headers {
			if h.Key == "" || h.Value == "" {
				continue
			}

			req.Header.Add(h.Key, h.Value)
		}

		start := time.Now()

		resp, err := client.Do(req)
		if err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return
		}
		defer resp.Body.Close()

		d := time.Since(start)
		r.Duration = d.Milliseconds()

		r.Status = resp.Status
		r.StatusCode = resp.StatusCode
		r.ContentType = resp.Header.Get("Content-Type")

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return
		}

		r.Size = int64(len(bodyBytes))

		sBody := string(bodyBytes)

		// deterministic formatting (json & html)
		if strings.HasPrefix(r.ContentType, "application/json") {
			var v any
			if err := json.Unmarshal(bodyBytes, &v); err != nil {
				App.Event.Emit("errMsg", fmt.Sprint(err))
				return
			}

			j, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				App.Event.Emit("errMsg", fmt.Sprint(err))
				return
			}

			sBody = string(j)

		} else if strings.HasPrefix(r.ContentType, "text/html") {
			sBody = gohtml.Format(sBody)
		}

		// defaults to string if content-type isnt supported
		r.Body = sBody
		s.Response = r // store within the scoop

		// emit event with Scoop containing the response
		App.Event.Emit("respMsg", s)
	}()
}

func (b *ScoopService) CreateCollection(c Collection) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) CreateScoop(c Collection, s Scoop) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// add scoop to current collection (no response)
	c.Scoops = append(c.Scoops, s)

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) DeleteScoop(c Collection, s Scoop) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// remove scoop from current collection
	c.Scoops = slices.DeleteFunc(c.Scoops, func(e Scoop) bool {
		return e.ID == s.ID
	})

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) OpenCollections() ([]Collection, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return nil, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	coll, err := os.ReadDir(scoopDir)
	if err != nil {
		// App.Event.Emit("errMsg", fmt.Sprint(err))
		// dont need to emit error here, frontend will just display no collections
		return nil, err
	}

	var availCollections []Collection
	for _, c := range coll {
		path := filepath.Join(scoopDir, c.Name())
		ext := filepath.Ext(path)

		if ext != ".json" {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return nil, err
		}

		var tempColl Collection
		if err := json.Unmarshal(content, &tempColl); err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return nil, err
		}

		availCollections = append(availCollections, tempColl)
	}

	return availCollections, nil
}

func (b *ScoopService) SaveCollection(c Collection) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) DeleteCollection(c Collection) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	// delete the entire file since each file defines a collection
	if err := os.Remove(path); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) SaveScoop(s Scoop, c Collection) (bool, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	scoopDir := filepath.Join(base, "Scoop", "Collections")

	// ensure /Scoop/Collections/ is created in UserConfigDir
	if err := os.MkdirAll(scoopDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	for i, scoop := range c.Scoops {
		if scoop.ID == s.ID {
			c.Scoops[i] = s
		}
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.ID))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *ScoopService) OpenDNSOverrides() (allOv []DNSOverride, ovDir string, err error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, "", err
	}

	dnsDir := filepath.Join(base, "Scoop", "DNS")

	// ensure /Scoop/DNS is created in UserConfigDir
	if err := os.MkdirAll(dnsDir, 0o755); err != nil {
		return nil, "", err
	}

	overrides := filepath.Join(dnsDir, "overrides.json")

	o, err := os.ReadFile(overrides)
	if err != nil {
		// on error first check the file exists
		if os.IsNotExist(err) {

			// if the file doesnt exist, create it
			if err := os.WriteFile(overrides, nil, 0o644); err != nil {
				return nil, "", err
			}

			// at this point file is empty so set 'o' to empty slice of bytes
			o = []byte{}

		} else {
			return nil, "", err
		}
	}

	if len(o) == 0 {
		return []DNSOverride{}, overrides, nil
	}

	if err := json.Unmarshal(o, &allOv); err != nil {
		return nil, "", err
	}

	return allOv, overrides, nil
}

func (b *ScoopService) CheckDNSOverride(s Scoop) (realURL string, err error) {
	allOV, _, err := b.OpenDNSOverrides()
	if err != nil {
		return "", err
	}

	for _, ov := range allOV {
		if strings.Contains(s.Request.URL, fmt.Sprintf("//%s", ov.Variable)) {
			realURL = strings.Replace(s.Request.URL, ov.Variable, ov.IPV4, 1)
			return
		}
	}

	return realURL, nil
}

func (b *ScoopService) CreateDNSOverride(newOv DNSOverride) (bool, error) {
	allOv, ovDir, err := b.OpenDNSOverrides()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	allOv = append(allOv, newOv)

	j, err := json.MarshalIndent(allOv, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	if err := os.WriteFile(ovDir, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

// the function below is for receiving ALL OVs from the sync server,
// do not want to append them to the existing ovs, need to overwrite the entire file

func (b *ScoopService) OverwriteDNSOverride(ov []DNSOverride) (bool, error) {
	_, ovDir, err := b.OpenDNSOverrides()
	if err != nil {
		return false, err
	}

	j, err := json.MarshalIndent(ov, "", "  ")
	if err != nil {
		return false, err
	}

	if err := os.WriteFile(ovDir, j, 0o644); err != nil {
		return false, err
	}

	return true, nil
}

func (b *ScoopService) GenerateCurlCommand(s Scoop) (string, error) {
	// add query params to url
	AddQueryParams(&s)

	// check for use of DNS Overrides
	realURL, err := b.CheckDNSOverride(s)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return "", err
	}

	// incase no matches
	if realURL == "" {
		realURL = s.Request.URL
	}

	req, err := http.NewRequest(string(s.Request.Method), realURL, nil)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return "", err
	}

	// add headers to request
	for _, h := range s.Request.Headers {
		if h.Key == "" || h.Value == "" {
			continue
		}

		req.Header.Add(h.Key, h.Value)
	}

	command, err := http2curl.GetCurlCommand(req)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return "", err
	}

	return command.String(), nil
}
