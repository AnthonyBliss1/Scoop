package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type ServerPayload struct {
	Collections []Collection  `json:"collections"`
	DNS         []DNSOverride `json:"dns"`
}

type SyncServer struct {
	ScoopService
}

var client = &http.Client{}

func (b *SyncServer) SetSyncServer(s Server) (ok bool, err error) {
	base, err := os.UserConfigDir()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	sServerDir := filepath.Join(base, "Scoop", "Sync-Server")

	// ensure /Scoop/Sync-Server is created in UserConfigDir
	if err := os.MkdirAll(sServerDir, 0o755); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	serverPath := filepath.Join(sServerDir, "server.json")

	j, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	if err := os.WriteFile(serverPath, j, 0o644); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *SyncServer) OpenSyncServer() (s Server, err error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return s, err
	}

	sServerDir := filepath.Join(base, "Scoop", "Sync-Server")

	// ensure /Scoop/Sync-Server is created in UserConfigDir
	if err := os.MkdirAll(sServerDir, 0o755); err != nil {
		return s, err
	}

	serverPath := filepath.Join(sServerDir, "server.json")

	data, err := os.ReadFile(serverPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Write Emtpy file if server file doesnt exist
			if err := os.WriteFile(serverPath, []byte{}, 0o644); err != nil {
				return s, err
			} else {
				// after writing the empty file, return empty Server Object
				return s, nil
			}
		}
		return s, err
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return s, err
	}

	return s, nil
}

func (b *SyncServer) SendToServer(s Server) (ok bool, err error) {
	if s.URL == "" {
		return false, errors.New("No server URL found")
	}

	// grab collections
	colls, err := b.OpenCollections()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// grab dnsOverrides
	dns, _, err := b.OpenDNSOverrides()
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// craft into a payload object
	payload := ServerPayload{Collections: colls, DNS: dns}

	// marshal to bytes (for the req body)
	data, err := json.Marshal(payload)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}
	reqBody := bytes.NewReader(data)

	// build req with /upload path in url
	req, err := http.NewRequest("POST", s.URL+"/upload", reqBody)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// make sure the server knows the type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.Key)

	resp, err := client.Do(req)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("POST to server failed: %q", bodyBytes)
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	return true, nil
}

func (b *SyncServer) GetFromServer(s Server) (ok bool, err error) {
	if s.URL == "" {
		return false, errors.New("No server URL found")
	}

	req, err := http.NewRequest("GET", s.URL+"/sync", nil)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	req.Header.Add("X-API-Key", s.Key)

	resp, err := client.Do(req)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("GET from server failed: %q", bodyBytes)
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	j, err := io.ReadAll(resp.Body)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	var payload ServerPayload
	if err := json.Unmarshal(j, &payload); err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	// save the payload data, starting with Collections
	// TODO: this does not override collections that were not pushed to sync server
	// consider the action and determine how I want this to operate
	for _, c := range payload.Collections {
		if _, err := b.CreateCollection(c); err != nil {
			App.Event.Emit("errMsg", fmt.Sprint(err))
			return false, err
		}
	}

	// save DNSOverrides
	if _, err := b.OverwriteDNSOverride(payload.DNS); err != nil {
		App.Event.Emit("errMsg", fmt.Sprintf("OverwriteDNS: %q", err))
		return false, err
	}

	return true, nil
}

func (b *SyncServer) CheckServerHealth(s Server) (ok bool, err error) {
	if s.URL == "" {
		err = errors.New("No server URL found")
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	req, err := http.NewRequest("GET", s.URL+"/health", nil)
	if err != nil {
		App.Event.Emit("errMsg", fmt.Sprint(err))
		return false, err
	}

	req.Header.Add("X-API-Key", s.Key)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return resp.StatusCode == http.StatusOK, nil
}
