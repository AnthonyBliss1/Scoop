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
	Name    string `json:"name"`
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
	ContentType string `json:"content_type"`
}

type Scoop struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Collection struct {
	Name   string  `json:"name"`
	Scoops []Scoop `json:"scoops"`
}

type Backend struct {
	context context.Context
}

// Initializes the Scoop model by attaching method, url, headers, query params, and body (body soon)

func (b *Backend) ModelIntializer(name string, method Method, reqURL string, headers []KV, qParams []KV) (Request, error) {
	var r Request

	r.Name = name
	r.Method = method
	r.URL = reqURL
	r.Headers = headers
	r.QParams = qParams

	return r, nil
}

func (b *Backend) AddQueryParams(s Scoop) error {
	u, err := url.Parse(s.Request.URL)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return err
	}
	query := url.Values{}

	for _, param := range s.Request.QParams {
		query.Add(param.Key, param.Value)
	}

	u.RawQuery = query.Encode()

	s.Request.URL = u.String()

	return nil
}

func (b *Backend) SubmitRequest(s Scoop) {
	go func() {
		var r Response

		client := http.Client{}

		// add query params to url
		b.AddQueryParams(s)

		req, err := http.NewRequest(string(s.Request.Method), s.Request.URL, nil)
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

		sBody := string(bodyBytes)

		// deterministic formatting
		if strings.HasPrefix(r.ContentType, "application/json") {
			var v any
			if err := json.Unmarshal(bodyBytes, &v); err != nil {
				App.Event.Emit("errMsg", fmt.Sprint(err))
				return
			}

			b, err := json.MarshalIndent(v, "", "  ")
			if err != nil {
				App.Event.Emit("errMsg", fmt.Sprint(err))
				return
			}

			sBody = string(b)

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

func (b *Backend) CreateCollection(c Collection) (bool, error) {
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

func (b *Backend) CreateRequest(c Collection, r Request) (bool, error) {
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

	// add request to current collection (no response)
	c.Scoops = append(c.Scoops, Scoop{Request: r})

	j, err := json.MarshalIndent(c.Scoops, "", "  ")
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

func (b *Backend) OpenCollections() ([]Collection, error) {
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

// TODO
// nearly identical to the CreateCollection function
// should make one function for create and save (WriteCollection)

func (b *Backend) SaveCollection(c Collection) (bool, error) {
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

func (b *Backend) SaveRequest(r Request, c Collection) (bool, error) {
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
		if scoop.Request.Name == r.Name {
			scoop.Request = r
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
