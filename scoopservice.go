package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type Backend struct {
	ctx context.Context
}

// Initializes the Scoop model by attaching method, url, headers, query params, and body (body soon)

func (a *Backend) ModelIntializer(method Method, reqURL string, headers []KV, qParams []KV) (*Scoop, error) {
	var scoop Scoop

	scoop.Request.Method = method
	scoop.Request.URL = reqURL
	scoop.Request.Headers = headers
	scoop.Request.QParams = qParams

	return &scoop, nil
}

// Simple func to submit request and store response back to the struct

func (a *Backend) SubmitRequest(s *Scoop) {
	go func() {
		var r Response

		client := http.Client{}

		// add query params to url
		a.AddQueryParams(s)

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

		App.Event.Emit("respMsg", s)
	}()
}

func (a *Backend) AddQueryParams(s *Scoop) error {
	u, err := url.Parse(s.Request.URL)
	if err != nil {
		App.Event.Emit("errMsg", err)
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
