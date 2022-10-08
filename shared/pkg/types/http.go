package types

type HTTPInfo struct {
	Version   string `json:"version,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	Origin    string `json:"origin,omitempty"`
	Host      string `json:"host,omitempty"`
}
