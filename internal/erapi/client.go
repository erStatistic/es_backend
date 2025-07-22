package erapi

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
	apiKey     string
}

func NewClient(timeout time.Duration, key string) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		apiKey: key,
	}
}
