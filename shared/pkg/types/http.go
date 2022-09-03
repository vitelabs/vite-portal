package types

type HttpInfo struct {
	Version   string `json:"version"`
	UserAgent string `json:"userAgent"`
	Origin    string `json:"origin"`
	Host      string `json:"host"`
	Auth      string `json:"auth"`
}
