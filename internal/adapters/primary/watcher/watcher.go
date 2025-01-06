package watcher

import (
	"context"

	"github.com/awlsring/dynamic-ip-watcher/internal/ports/service"
)

type Watcher struct {
	addressService service.Address
}

func New(addressService service.Address) *Watcher {
	return &Watcher{
		addressService: addressService,
	}
}

func (w *Watcher) Run(ctx context.Context) error {
	return w.addressService.DetectAndHandleAddressChange(ctx)
}
