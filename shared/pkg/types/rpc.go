package types

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}