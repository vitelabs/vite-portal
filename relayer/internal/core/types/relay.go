package types

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// A read/write request to be relayed
type Relay struct {
	Config  RelayConfig
	Payload Payload
}

type RelayConfig struct {
	Debug            bool
	UserAgent        string
	RpcNodeTimeout   int64
	SortJsonResponse bool
}

// The data being sent to the external node
type Payload struct {
	Data    string              `json:"data"`              // The actual data string to be sent
	Method  string              `json:"method"`            // The HTTP CRUD method
	Path    string              `json:"path"`              // The REST path
	Headers map[string][]string `json:"headers,omitempty"` // The HTTP headers
}

// The response to a relay request
type RelayResponse struct {
	Response string
}

// Execute attempts to do a request on the specified node
func (r Relay) Execute() (string, roottypes.Error) {
	// TODO: get node HTTP RPC url
	nodeHttpRpcUrl := "http://127.0.0.1:23456"
	url := strings.Trim(nodeHttpRpcUrl, "/")
	if len(r.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(r.Payload.Path, "/")
	}
	res, err := executeHttpRequest(&r.Config, r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if err != nil {
		// TODO: track metrics
		return res, NewError(DefaultCodeNamespace, CodeHttpExecutionError, err)
	}
	return res, nil
}

// executeHttpRequest takes in the raw json string and forwards it to the RPC endpoint
func executeHttpRequest(cfg *RelayConfig, payload, url, method string, headers map[string][]string) (string, error) {
	// generate the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if cfg.UserAgent != "" {
		req.Header.Set("User-Agent", cfg.UserAgent)
	}
	// TODO: set basic auth instead of IP whitelisting
	// req.SetBasicAuth(username, password)
	if len(headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range headers {
			for _, s := range v {
				req.Header.Set(k, s)
			}
		}
	}
	// execute the request
	resp, err := (&http.Client{Timeout: time.Duration(cfg.RpcNodeTimeout) * time.Millisecond}).Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res := string(body)
	if cfg.SortJsonResponse {
		res = sortJsonResponse(res)
	}
	return res, nil
}

// sortJsonResponse sorts json from a relay response
func sortJsonResponse(r string) string {
	var rawJSON map[string]interface{}
	// unmarshal into json
	if err := json.Unmarshal([]byte(r), &rawJSON); err != nil {
		return r
	}
	// marshal into json
	res, err := json.Marshal(rawJSON)
	if err != nil {
		return r
	}
	return string(res)
}
