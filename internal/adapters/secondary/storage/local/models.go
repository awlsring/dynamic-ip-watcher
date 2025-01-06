package local_storage

import (
	"net"
	"time"
)

type LastKnownIPAddressData struct {
	IPAddress net.IP    `json:"ip_address"`
	CheckedAt time.Time `json:"checked_at"`
}
