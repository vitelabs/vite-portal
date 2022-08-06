package store

import (
	"bytes"
	"net/http"
	"time"

	"github.com/vitelabs/vite-portal/internal/core/types"
	"github.com/vitelabs/vite-portal/internal/util/jsonutil"
)

type HttpCollector struct {
	url string
	userAgent string
	timeout int64
}

func NewHttpCollector(url, userAgent string, timeout int64) *HttpCollector {
	return &HttpCollector{
		url: url,
		userAgent: userAgent,
		timeout: timeout,
	}
}

func (c *HttpCollector) DispatchRelayResult(result types.RelayResult) error {
	data, err := jsonutil.ToByte(result)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err1 := (&http.Client{Timeout: time.Duration(c.timeout) * time.Millisecond}).Do(req)
	if err1 != nil {
		return err1
	}
	return nil
}