package ipapi

import "net/http"

type Option func(*IPAPIClient)

func WithHTTPClient(client *http.Client) Option {
	return func(r *IPAPIClient) {
		r.httpClient = client
	}
}

func WithEndpoint(endpoint string) Option {
	return func(r *IPAPIClient) {
		r.endpoint = endpoint
	}
}
