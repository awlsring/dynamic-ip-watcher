package service

import "context"

type Address interface {
	DetectAndHandleAddressChange(context.Context) error
}
