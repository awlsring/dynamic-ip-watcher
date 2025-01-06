package gateway

import (
	"context"
	"net"
)

type IPRetriever interface {
	GetPublicIPv4(context.Context) (net.IP, error)
}
