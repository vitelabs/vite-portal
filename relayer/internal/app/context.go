package app

import (
	nodeinterfaces "github.com/vitelabs/vite-portal/internal/node/interfaces"
	nodestore "github.com/vitelabs/vite-portal/internal/node/store"
)

type Context struct {
	nodeStore nodeinterfaces.StoreI
}

func NewContext() *Context {
	c := &Context{
		nodeStore: nodestore.NewKvStore(),
	}
	return c
}

func InitContext() (*Context, error) {
	c := NewContext()
	return c, nil
}