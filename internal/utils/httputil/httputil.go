// package httputil is a wrapper for net/http Client.
// httputil is WIP.
package httputil

import (
	"net"
	"net/http"
	"time"
)

// DefaultTimeout contains dial and requests default timeout values.
const DefaultTimeout = 10 * time.Second

// Config contains configuration to net/http client wrapper.
type Config struct {
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Delay   time.Duration `json:"delay" yaml:"yaml"`
}

// A Client is wrapper for net/http with dial and requests timeouts.
type Client struct {
	Config *Config
	http.Client
}

// NewClient creates new Client instance.
func NewClient(cfg *Config) *Client {
	c := Client{
		Config: cfg,
		Client: http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Dial:                (&net.Dialer{Timeout: cfg.Timeout}).Dial,
				TLSHandshakeTimeout: cfg.Timeout,
			},
		},
	}

	return &c
}

// NewClient creates new Client instance with default timeout.
func NewDefaultClient() *Client {
	return NewClient(&Config{Timeout: DefaultTimeout})
}
