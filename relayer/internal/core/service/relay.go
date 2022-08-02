package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/cryptoutil"
)

// HandleRelay handles a read/write request to one or multiple nodes
func (s *Service) HandleRelay(r coretypes.Relay) (*coretypes.RelayResponse, roottypes.Error) {
	// TODO: get node HTTP RPC url
	nodeHttpRpcUrl := "http://127.0.0.1:23456"
	url := strings.Trim(nodeHttpRpcUrl, "/")
	if len(r.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(r.Payload.Path, "/")
	}
	response, err := s.executeHttpRequest(r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("could not execute relay")
		return nil, coretypes.NewError(coretypes.DefaultCodeNamespace, coretypes.CodeHttpExecutionError, err)
	}
	res := &coretypes.RelayResponse{
		Response: response,
	}
	// TODO: track relay time and add to metrics
	return res, nil
}

// getConsensusNodes returns random nodes used for consensus
func (s *Service) getConsensusNodes(r coretypes.Relay) ([]nodetypes.Node, roottypes.Error) {
	header := coretypes.NewSessionHeader(r.ClientIp, r.Chain)
	session, err := s.HandleSession(header)
	if err != nil {
		return nil, err
	}
	if s.config.ConsensusNodeCount >= len(session.Nodes) {
		return session.Nodes, nil
	}
	sessionNodes := make([]nodetypes.Node, s.config.ConsensusNodeCount)
	rnd := cryptoutil.UniqueRandomInt(len(session.Nodes), s.config.ConsensusNodeCount)
	index := 0
	for _, v := range rnd {
		sessionNodes[index] = session.Nodes[v]
		index++
	}
	return sessionNodes, nil
}

// executeHttpRequest takes in the raw json string and forwards it to the RPC endpoint
func (s *Service) executeHttpRequest(payload, url, method string, headers map[string][]string) (string, error) {
	// generate the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	if s.config.UserAgent != "" {
		req.Header.Set("User-Agent", s.config.UserAgent)
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
	resp, err := (&http.Client{Timeout: time.Duration(s.config.RpcNodeTimeout) * time.Millisecond}).Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res := string(body)
	if s.config.SortJsonResponse {
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