package types

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/vitelabs/vite-portal/internal/logger"
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	"github.com/vitelabs/vite-portal/internal/types"
	"github.com/vitelabs/vite-portal/internal/util/cryptoutil"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
	"github.com/vitelabs/vite-portal/internal/util/mathutil"
)

// Session randomly groups one client with a set of nodes and is valid for a limited timeframe
type Session struct {
	Key       string           `json:"key"`
	Timestamp int64            `json:"timestamp"`
	Header    SessionHeader    `json:"header"`
	Nodes     []nodetypes.Node `json:"nodes"`
}

// SessionHeader defines the header for session information
type SessionHeader struct {
	IpAddress string `json:"ipAddress"`
	Chain     string `json:"chain"`
}

// NewSession creates a new session
func NewSession(s nodeinterfaces.ServiceI, header SessionHeader, nodeCount int) (Session, types.Error) {
	sessionNodes, err := NewSessionNodes(s, header.Chain, nodeCount)
	if err != nil {
		return Session{}, err
	}
	return Session{
		Key:       header.HashString(),
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
		return NewError(DefaultCodeNamespace, CodeInvalidChainError, errors.New("empty"))
	}
	// verify the ip address
	if sh.IpAddress == "" {
		return NewError(DefaultCodeNamespace, CodeInvalidIpAddressError, errors.New("empty"))
	}
	return nil
}

// Hash generates the cryptographic hash representation of the session header
func (sh SessionHeader) Hash() []byte {
	res := md5.Sum(sh.Bytes())
	return res[:]
}

// HashString generates the hex string representation of the cryptographic hash
func (sh SessionHeader) HashString() string {
	return hex.EncodeToString(sh.Hash())
}

// Bytes generates the bytes representation of the session header
func (sh SessionHeader) Bytes() []byte {
	res, err := jsonutil.ToByte(sh)
	if err != nil {
		logger.Logger().Fatal().Err(err).Msg("an error occurred trying to convert the session key into bytes")
	}
	return res
}
