package iputil

import (
	"fmt"
	"net"
	"net/http"
)

const (
	IpHeader = "X-Real-IP"
)

func IsTrusted(ip string, trustedSubnet string) (bool, error) {
	if trustedSubnet == "" {
		return false, nil
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, fmt.Errorf("invalid ip address: %s", ip)
	}

	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false, fmt.Errorf("invalid trusted subnet: %s, err: %s", trustedSubnet, err.Error())
	}

	return ipNet.Contains(parsedIP), nil
}

func IpFromRequest(request *http.Request) string {
	return request.Header.Get(IpHeader)
}
