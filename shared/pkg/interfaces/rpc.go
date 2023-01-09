package interfaces

import "context"

type RpcClientI interface {
	GetID() uint32
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Notify(ctx context.Context, method string, args ...interface{}) error 
}

type RpcResponseI interface {
	GetError() error
}