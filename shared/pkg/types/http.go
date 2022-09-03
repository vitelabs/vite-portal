package types

type HttpInfo struct {
	Version   string `json:"id"`
	UserAgent string `json:"userAgent"`
	Origin    string `json:"origin"`
	Host      string `json:"host"`
	Auth      string `json:"auth"`
}
