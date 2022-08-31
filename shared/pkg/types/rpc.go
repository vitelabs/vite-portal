package types

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RpcAppInfoResponse struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Name    string `json:"name"`
}
