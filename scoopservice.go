package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	Name     string   `json:"name"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Collection struct {
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

// Initializes the Scoop model by attaching method, url, headers, query params, and body (body soon)

func (b *ScoopService) ModelIntializer(method Method, reqURL string, headers []KV, qParams []KV) (Request, error) {
	var r Request

	r.Method = method
	r.URL = reqURL
	r.Headers = headers
	r.QParams = qParams

	return r, nil
}

func (b *ScoopService) AddQueryParams(s *Scoop) error {
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

func (b *ScoopService) SubmitRequest(s Scoop) {
	go func() {
		var r Response

		client := http.Client{}

		// add query params to url
		b.AddQueryParams(&s)

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

		req, err := http.NewRequest(string(s.Request.Method), realURL, nil)
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

// TODO: prevent duplicate Collection names

func (b *ScoopService) CreateCollection(c Collection) (bool, error) {
	if strings.ContainsAny(c.Name, `/\`) {
		err := errors.New("collection name cannot contain slashes")
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

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

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.Name))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

// TODO: prevent duplicate Scoop names

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

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.Name))
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

// TODO:
// nearly identical to the CreateCollection function
// should make one function for create and save (WriteCollection)

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

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.Name))
	path := filepath.Join(scoopDir, colFile)

	if err := os.WriteFile(path, j, 0o644); err != nil {
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

	for _, scoop := range c.Scoops {
		if scoop.Name == s.Name {
			scoop = s
		}
	}

	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	colFile := fmt.Sprintf("%s.json", strings.TrimSpace(c.Name))
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

func (b *ScoopService) GenerateCurlCommand(s Scoop) (string, error) {
	// add query params to url
	b.AddQueryParams(&s)

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
