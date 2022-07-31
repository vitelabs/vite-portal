package service

import (
	"errors"

	coretypes "github.com/vitelabs/vite-portal/internal/core/types"
	nodetypes "github.com/vitelabs/vite-portal/internal/node/types"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// HandleSession handles the session information for a client request
func (s *Service) HandleSession(header coretypes.SessionHeader) (*coretypes.Session, roottypes.Error) {
	err := header.ValidateHeader()
	if err != nil {
		return nil, err
	}
	nodeCount := s.nodeService.GetNodeCount(header.Chain)
	if nodeCount == 0 {
		return nil, coretypes.NewError(coretypes.DefaultCodeNamespace, coretypes.CodeInvalidChainError, errors.New("no nodes"))
	}

	session, err := s.getSession(header)
	if err != nil {
		return nil, err
	}

	// check if nodes have been deleted or updated since last time using the session
	if s.haveNodesChanged(session) {
		actualNodes := s.getActualNodes(session)
		session.Nodes = actualNodes
		s.cache.SetSession(session)
	}

	minNodeCount := s.config.ConsensusNodeCount
	// make sure session has sufficient nodes
	if nodeCount > minNodeCount && minNodeCount > len(session.Nodes) || len(session.Nodes) == 0 {
		// delete current session and create new
		s.cache.DeleteSession(header)
		session, err = s.getSession(header)
		if err != nil {
			return nil, err
		}
	}

	return &session, nil
}

func (s *Service) haveNodesChanged(session coretypes.Session) bool {
	hasDeletedNodes := s.nodeService.LastActivityTimestamp(session.Header.Chain, nodetypes.Delete) > session.Timestamp
	hasUpdatedNodes := s.nodeService.LastActivityTimestamp(session.Header.Chain, nodetypes.Put) > session.Timestamp
	return hasDeletedNodes || hasUpdatedNodes
}

func (s *Service) getSession(header coretypes.SessionHeader) (coretypes.Session, roottypes.Error) {
	// check cache
	session, found := s.cache.GetSession(header, s.config.MaxSessionDuration)
	if !found {
		// create new session
		newSession, err := coretypes.NewSession(s.nodeService, header, s.config.SessionNodeCount)
		if err != nil {
			return coretypes.Session{}, err
		}
		// add to cache
		s.cache.SetSession(newSession)
		session = newSession
	}
	return session, nil
}

func (s *Service) getActualNodes(session coretypes.Session) []nodetypes.Node {
	var actualNodes []nodetypes.Node
	for _, v := range session.Nodes {
		n, found := s.nodeService.GetNode(v.Id)
		if !found || n.Chain != session.Header.Chain {
			continue
		}
		actualNodes = append(actualNodes, n)
	}
	if len(actualNodes) == 0 {
		return []nodetypes.Node{}
	}
	return actualNodes
}