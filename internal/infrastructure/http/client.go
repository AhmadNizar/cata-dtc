package http

import (
	"net/http"
	"time"
)

type Config struct {
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
}

func NewHTTPClient(config Config) *http.Client {
	return &http.Client{
		Timeout: config.Timeout,
	}
}