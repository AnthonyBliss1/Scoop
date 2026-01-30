package main

import (
	"context"
	"io"
	"net/http"
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
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Headers    []KV   `json:"headers"`
	Body       string `json:"body"`
}

type Scoop struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type App struct {
	ctx context.Context
}

// Needs to take Headers, Body, Params - just doing headers for now

func (a *App) ModelIntializer(method Method, reqURL string, headers []KV) (*Scoop, error) {
	var scoop Scoop

	scoop.Request.Method = method
	scoop.Request.URL = reqURL
	scoop.Request.Headers = headers

	return &scoop, nil
}

// Simple func to submit request and store response back to the struct

func (a *App) SubmitRequest(s *Scoop) (Response, error) {
	var r Response

	client := http.Client{}

	req, err := http.NewRequest(string(s.Request.Method), s.Request.URL, nil)
	if err != nil {
		return Response{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	r.Status = resp.Status
	r.StatusCode = resp.StatusCode

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	r.Body = string(bodyBytes)

	return r, nil
}
