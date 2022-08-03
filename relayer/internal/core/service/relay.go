package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/logger"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/cryptoutil"
)

// HandleRelay handles a read/write request to one or multiple nodes
func (s *Service) HandleRelay(r coretypes.Relay) (*coretypes.RelayResponse, roottypes.Error) {
	nodes, err := s.getConsensusNodes(r)
	if err != nil {
		return nil, err
	}
	err1 := errors.New("relay timed out")
	c := make(chan string, len(nodes))
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(s.config.RpcNodeTimeout))
	defer cancelFn()
	var processed uint32
	// Relay to multiple nodes and return the fastest response
	for _, n := range nodes {
		go func(n nodetypes.Node) {
			response, err := s.executeRelay(ctx, n.RpcHttpUrl, r)
			atomic.AddUint32(&processed, 1)
			if err != nil {
				if int(processed) >= len(nodes) {
					err1 = errors.New("relay failed")
					// No more nodes left -> cancel the relay
					cancelFn()
				}
				return
			}
			c <- response
		}(n)
	}
	for {
		select {
		case response := <- c:
			res := &coretypes.RelayResponse{
				Response: response,
			}
			// TODO: track relay time and add to metrics
			return res, nil
		case <-ctx.Done():
			logger.Logger().Error().Err(err1).Msg("relay cancelled")
			return &coretypes.RelayResponse{}, coretypes.NewError(coretypes.DefaultCodeNamespace, coretypes.CodeHttpExecutionError, err1)
		}
	}
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

func (s *Service) executeRelay(ctx context.Context, nodeHttpRpcUrl string, r coretypes.Relay) (string, roottypes.Error) {
	url := strings.Trim(nodeHttpRpcUrl, "/")
	if len(r.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(r.Payload.Path, "/")
	}
	response, err := s.executeHttpRequest(ctx, r.Payload.Data, url, r.Payload.Method, r.Payload.Headers)
	if err != nil {
		logger.Logger().Error().Err(err).Msg("could not execute relay")
		return "", coretypes.NewError(coretypes.DefaultCodeNamespace, coretypes.CodeHttpExecutionError, err)
	}
	return response, nil
}

// executeHttpRequest takes in the raw json string and forwards it to the RPC endpoint
func (s *Service) executeHttpRequest(ctx context.Context, payload, url, method string, headers map[string][]string) (string, error) {
	// generate the request
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer([]byte(payload)))
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
