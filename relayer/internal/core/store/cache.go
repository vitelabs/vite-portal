package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/vitelabs/vite-portal/relayer/internal/core/types"
	roottypes "github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/logger"
)

type CacheStore struct {
	sessionCache roottypes.Cache
	lock         sync.RWMutex
}

func NewCacheStore(capacity int) *CacheStore {
	s := &CacheStore{
		sessionCache: *roottypes.NewCache(capacity),
		lock:         sync.RWMutex{},
	}
	return s
}

// GetSession returns a session (value) from the cache using a header (key)
func (s *CacheStore) GetSession(header types.SessionHeader, maxDuration int64) (session types.Session, found bool) {
	key := getSessionKey(header)
	val, found := s.sessionCache.Get(key)
	if !found {
		return types.Session{}, false
	}
	sess, ok := val.(types.Session)
	if !ok {
		logger.Logger().Error().Msg(fmt.Sprintf("could not unmarshal into session from cache with header %v", header))
		return types.Session{}, false
	}
	// check if expired
	if time.Now().UnixMilli()-maxDuration > sess.Timestamp {
		return types.Session{}, false
	}
	return sess, true
}

// SetSession sets a session (value) in the cache using the header (key)
func (s *CacheStore) SetSession(session types.Session) {
	key := getSessionKey(session.Header)
	s.sessionCache.Add(key, session)
}

// DeleteSession deletes a session (value) from the cache
func (s *CacheStore) DeleteSession(header types.SessionHeader) {
	key := getSessionKey(header)
	s.sessionCache.Remove(key)
}

// ClearSessions clears all sessions from the cache
func (s *CacheStore) ClearSessions() {
	s.sessionCache.Purge()
}

func getSessionKey(header types.SessionHeader) string {
	return header.HashString()
}
