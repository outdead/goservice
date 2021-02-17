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
	Timeout             time.Duration `json:"timeout" yaml:"timeout"`
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`
	Delay               time.Duration `json:"delay" yaml:"yaml"`
}

// A Client is wrapper for net/http with dial and requests timeouts.
type Client struct {
	Config *Config
	http.Client
}

// NewClient creates new Client instance.
func NewClient(cfg *Config) *Client {
	tlsHandshakeTimeout := cfg.Timeout
	if cfg.TLSHandshakeTimeout != 0 {
		tlsHandshakeTimeout = cfg.TLSHandshakeTimeout
	}

	c := Client{
		Config: cfg,
		Client: http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Dial:                (&net.Dialer{Timeout: cfg.Timeout}).Dial,
				TLSHandshakeTimeout: tlsHandshakeTimeout,
			},
		},
	}

	return &c
}

// NewClient creates new Client instance with default timeout.
func NewDefaultClient() *Client {
	return NewClient(&Config{Timeout: DefaultTimeout})
}
