package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/internal/logger"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/cryptoutil"
	"github.com/vitelabs/vite-portal/internal/util/mathutil"
)

// The response object used in dispatching
type DispatchResponse struct {
	Session Session `json:"session"`
}

type Session struct {
	Timestamp int64            `json:"timestamp"`
	Header    SessionHeader    `json:"header"`
	Nodes     []nodetypes.Node `json:"nodes"`
}

// SessionHeader defines the header for session information
type SessionHeader struct {
	IpAddress string `json:"ipAddress"`
	Chain     string `json:"chain"`
	Timestamp int64  `json:"timestamp"`
}

// NewSession creates a new session from seed data
func NewSession(s nodeinterfaces.ServiceI, header SessionHeader, nodeCount int) (Session, types.Error) {
	sessionNodes, err := NewSessionNodes(s, header.Chain, nodeCount)
	if err != nil {
		return Session{}, err
	}
	return Session{
		Timestamp: time.Now().UnixMilli(),
		Header:    header,
		Nodes:     sessionNodes,
	}, nil
}

// NewSessionNodes creates nodes for the session
func NewSessionNodes(s nodeinterfaces.ServiceI, chain string, nodeCount int) ([]nodetypes.Node, types.Error) {
	currentNodeCount := s.GetNodeCount(chain)
	if currentNodeCount <= 0 {
		return nil, NewBasicError(DefaultCodeNamespace, CodeInsufficientNodesError)
	}
	sessionNodeCount := mathutil.Min(nodeCount, currentNodeCount)
	sessionNodes := make([]nodetypes.Node, sessionNodeCount)
	r := cryptoutil.UniqueRandomInt(currentNodeCount, sessionNodeCount)
	index := 0
	for _, v := range r {
		node, found := s.GetNodeByIndex(chain, v)
		if !found {
			logger.Logger().Info().Msg(fmt.Sprintf("inconsistent state when trying to get node by index for chain '%s'", chain))
			return NewSessionNodes(s, chain, nodeCount)
		}
		sessionNodes[index] = node
		index++
	}
	return sessionNodes, nil
}

// ValidateHeader validates the header of the session
func (sh SessionHeader) ValidateHeader() types.Error {
	// verify the chain
	if sh.Chain == "" {
		return NewError(DefaultCodeNamespace, CodeInvalidChain, errors.New("empty"))
	}
	// verify the ip address
	if sh.IpAddress == "" {
		return NewError(DefaultCodeNamespace, CodeInvalidIpAddress, errors.New("empty"))
	}
	// verify the timestamp
	if sh.Timestamp < 1 {
		return NewError(DefaultCodeNamespace, CodeInvalidTimestamp, errors.New(fmt.Sprintf("%d", sh.Timestamp)))
	}
	return nil
}
