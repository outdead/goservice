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

// A Client is wrapper for net/http with dial and requests timeouts.
type Client struct {
	Dialer *Dialer
	http.Client
}

// NewClient creates new HTTP Client instance.
func NewClient(dialer *Dialer) *Client {
	tlsHandshakeTimeout := dialer.Timeout
	if dialer.TLSHandshakeTimeout != 0 {
		tlsHandshakeTimeout = dialer.TLSHandshakeTimeout
	}

	c := Client{
		Dialer: dialer,
		Client: http.Client{
			Timeout: dialer.Timeout,
			Transport: &http.Transport{
				Dial:                (&net.Dialer{Timeout: dialer.Timeout}).Dial,
				TLSHandshakeTimeout: tlsHandshakeTimeout,
			},
		},
	}

	return &c
}

// NewClient creates new Client instance with default timeout.
func NewDefaultClient() *Client {
	return NewClient(&Dialer{Timeout: DefaultTimeout})
}
