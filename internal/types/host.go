package types

// Package types provides common types used across the application.

import (
	"fmt"
	"net"
)

type Host string

func isHostname(s string) bool {
	// A lazy hostname check. You can replace this with real DNS validation if you hate yourself.
	return len(s) > 0 && len(s) <= 255
}

func NewHost(val string) (Host, error) {
	if net.ParseIP(val) == nil && !isHostname(val) {
		return "", fmt.Errorf("invalid bind address: '%s'", val)
	}
	return Host(val), nil
}
