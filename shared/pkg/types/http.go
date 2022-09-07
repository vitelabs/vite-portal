package types

type HTTPInfo struct {
	Version   string `json:"version"`
	UserAgent string `json:"userAgent"`
	Origin    string `json:"origin"`
	Host      string `json:"host"`
	Auth      string `json:"-"`
}
