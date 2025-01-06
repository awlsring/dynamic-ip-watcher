package gateway

import (
	"context"
	"errors"
	"net"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrMultipleRecordsFound = errors.New("multiple records found")
)

type DNSUpdater interface {
	RecordName() string
	GetRecordIpAddress(ctx context.Context) (net.IP, error)
	CreateRecordWithIpAddress(ctx context.Context, ip net.IP) error
	UpdateRecordIpAddress(ctx context.Context, ip net.IP) error
}
