// Package xpack holds the community implementation of the optional provider
// helpers used by agent: the shared HTTP transport for outbound alert delivery
// and the metadata describing the agent that sends an alert.
package xpack

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/1Panel-dev/1Panel/agent/app/dto"
)

// LoadRequestTransport returns the shared transport used for outbound requests.
func LoadRequestTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		IdleConnTimeout:       15 * time.Second,
	}
}

// GetAgentInfo returns nil in community edition, which always runs on the local
// agent. Callers fall back to the local node name and address.
func GetAgentInfo() (*dto.AgentInfo, error) {
	return nil, nil
}
