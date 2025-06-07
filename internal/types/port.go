package types

import "fmt"

type Port int

func NewPort(val int) (Port, error) {
	if val <= 0 || val > 65535 {
		return 0, fmt.Errorf("invalid port: %d", val)
	}
	return Port(val), nil
}
