package types

import (
	"fmt"

	"github.com/vitelabs/vite-portal/internal/logger"
	roottypes "github.com/vitelabs/vite-portal/internal/types"
)

// GetSession returns a session (value) from the cache using a header (key)
func GetSession(header SessionHeader) (session Session, found bool) {
	key := getSessionKey(header)
	val, found := roottypes.GlobalSessionCache.Get(key)
	if !found {
		return Session{}, found
	}
	session, found = val.(Session)
	if !found {
		logger.Logger().Error().Msg(fmt.Sprintf("could not unmarshal into session from cache with header %v", header))
	}
	return
}

// SetSession sets a session (value) in the cache using the header (key)
func SetSession(session Session) {
	key := getSessionKey(session.Header)
	roottypes.GlobalSessionCache.Add(key, session)
}

// DeleteSession deletes a session (value) from the cache
func DeleteSession(header SessionHeader) {
	key := getSessionKey(header)
	roottypes.GlobalSessionCache.Remove(key)
}

// ClearSessions clears all sessions from the cache
func ClearSessions() {
	if roottypes.GlobalSessionCache != nil {
		roottypes.GlobalSessionCache.Purge()
	}
}

func getSessionKey(header SessionHeader) string {
	return header.IpAddress
}