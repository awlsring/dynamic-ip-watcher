package ipapi_ip_retriever

import (
	"context"
	"net"

	ipapi "github.com/awlsring/dynamic-ip-watcher/internal/pkg/ip-api"
	"github.com/awlsring/dynamic-ip-watcher/internal/ports/gateway"
)

type IPRetrieverIPAPI struct {
	client ipapi.Client
}

func New(client ipapi.Client) gateway.IPRetriever {
	return &IPRetrieverIPAPI{
		client: client,
	}
}

func (r *IPRetrieverIPAPI) GetPublicIPv4(ctx context.Context) (net.IP, error) {
	response, err := r.client.GetPublicIP(ctx)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(response.Query), nil
}
