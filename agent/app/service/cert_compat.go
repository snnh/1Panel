package service

import (
	"crypto/x509"
	"net"
	"strings"

	"github.com/1Panel-dev/1Panel/agent/app/model"
	"github.com/1Panel-dev/1Panel/agent/constant"
	"golang.org/x/net/idna"
)

func websiteCertificateDomainMismatch(website model.Website, cert *x509.Certificate) bool {
	if website.Protocol != constant.ProtocolHTTPS || cert == nil {
		return false
	}
	for _, domain := range website.Domains {
		host := strings.TrimSpace(domain.Domain)
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		host = strings.TrimSuffix(strings.ToLower(strings.Trim(host, "[]")), ".")
		if host == "" {
			continue
		}
		if net.ParseIP(host) == nil {
			wildcard := strings.HasPrefix(host, "*.")
			normalized, err := idna.Lookup.ToASCII(strings.TrimPrefix(host, "*."))
			if err != nil {
				continue
			}
			host = normalized
			if wildcard {
				host = "*." + host
			}
		}
		if cert.VerifyHostname(host) != nil {
			return true
		}
	}
	return false
}
