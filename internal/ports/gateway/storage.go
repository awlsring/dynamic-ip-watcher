package gateway

import (
	"context"
	"net"
)

type Storage interface {
	SaveIPAddress(ctx context.Context, ip net.IP) error
	GetLastKnownIPAddress(ctx context.Context) (net.IP, error)
}
