package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	header := coretypes.NewSessionHeader(r.ClientIp, r.Chain)
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	nodes, err := s.getConsensusNodes(header)
	if err != nil {
		return nil, err
	}
	err1 := errors.New("relay timed out")
	c := make(chan string, len(nodes))
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(s.config.RpcNodeTimeout))
	_ = cancelFn // Ignore `lostcancel` warning (https://github.com/golang/go/issues/29587)
	var processed uint32
	responses := make([]coretypes.NodeResponse, len(nodes))
	// Relay to multiple nodes and return the fastest response
	for i, n := range nodes {
		go func(index int, n nodetypes.Node) {
			res := s.execute(ctx, n, r)
			responses[index] = res
			atomic.AddUint32(&processed, 1)
			if int(processed) >= len(nodes) {
				go s.dispatchRelayResult(r, header.HashString(), responses)
			}
			if res.Error != "" {
				if int(processed) >= len(nodes) {
					err1 = errors.New("relay failed")
					// No more nodes left -> cancel the relay
					cancelFn()
				}
				return
			}
			c <- res.Response
		}(i, n)
	}
	for {
		select {
		case response := <- c:
			res := &coretypes.RelayResponse{
				SessionKey: header.HashString(),
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
func (s *Service) getConsensusNodes(header coretypes.SessionHeader) ([]nodetypes.Node, roottypes.Error) {
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

// execute forwards the relay request and composes the node response
func (s *Service) execute(ctx context.Context, n nodetypes.Node, r coretypes.Relay) coretypes.NodeResponse {
	url := strings.Trim(n.RpcHttpUrl, "/")
	if len(r.Payload.Path) > 0 {
		url = url + "/" + strings.Trim(r.Payload.Path, "/")
	}
	startTime := time.Now()
	response, err := s.executeHttpRequest(ctx, url, r.Payload)
	result := coretypes.NodeResponse{
		NodeId: n.Id,
		ResponseTime: time.Since(startTime).Milliseconds(),
		Response: response,
	}
	if ctx.Err() != nil {
		if ctx.Err().Error() == "context deadline exceeded" {
			result.DeadlineExceeded = true
		} else {
			result.Cancelled = true
		}
	}
	if err != nil {
		result.Error = err.Error()
	}
	return result
}

// executeHttpRequest forwards the relay request to the HTTP RPC endpoint
func (s *Service) executeHttpRequest(ctx context.Context, url string, payload coretypes.Payload) (string, error) {
	// generate the request
	req, err := http.NewRequestWithContext(ctx, payload.Method, url, bytes.NewBuffer([]byte(payload.Data)))
	if err != nil {
		return "", err
	}
	if s.config.UserAgent != "" {
		req.Header.Set("User-Agent", s.config.UserAgent)
	}
	// TODO: set basic auth instead of IP whitelisting
	// req.SetBasicAuth(username, password)
	if len(payload.Headers) == 0 {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for k, v := range payload.Headers {
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

// dispatchRelayResult dispatches the relay result to the orchestrator or directly to Kafka
func (s *Service) dispatchRelayResult(r coretypes.Relay, sessionKey string, responses []coretypes.NodeResponse) {
	result := coretypes.RelayResult{
		SessionKey: sessionKey,
		Relay: r,
		Responses: responses,
	}
	if s.config.Debug {
		logger.Logger().Debug().Str("result", fmt.Sprintf("%#v", result)).Msg("relay result")
	}
	// TODO: send to Kafka
	if s.httpCollector != nil {
		err := s.httpCollector.DispatchRelayResult(result)
		if err != nil {
			logger.Logger().Error().Err(err).Msg("HttpCollector.DispatchRelayResult failed")
		}
	}
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
