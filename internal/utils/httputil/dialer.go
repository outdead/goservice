package httputil

import "time"

// Dialer contains configuration to net/http client wrapper.
type Dialer struct {
	Timeout             time.Duration `json:"timeout" yaml:"timeout"`
	TLSHandshakeTimeout time.Duration `json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`
	Delay               time.Duration `json:"delay" yaml:"yaml"`
}
