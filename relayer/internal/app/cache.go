package app

import (
	"github.com/vitelabs/vite-portal/internal/types"
)

func InitCache(capacity int) error {
	types.GlobalSessionCache = types.NewCache(capacity)
	return nil
}