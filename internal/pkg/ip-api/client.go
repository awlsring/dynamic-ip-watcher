package ipapi

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	IPAPIEndpoint = "http://ip-api.com/json/"
)

type Client interface {
	GetPublicIP(ctx context.Context) (*IPQueryResponse, error)
	QueryIPAddress(ctx context.Context, ip string) (*IPQueryResponse, error)
}

type IPAPIClient struct {
	httpClient *http.Client
	endpoint   string
}

func New(opts ...Option) *IPAPIClient {
	retriever := &IPAPIClient{
		httpClient: http.DefaultClient,
		endpoint:   IPAPIEndpoint,
	}

	for _, opt := range opts {
		opt(retriever)
	}

	return retriever
}

type IPQueryResponse struct {
	Status        string  `json:"status"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"lat"`
	Longitude     float64 `json:"lon"`
	Offset        int     `json:"offset"`
	Currency      string  `json:"currency"`
	Timezone      string  `json:"timezone"`
	ISP           string  `json:"isp"`
	ORG           string  `json:"org"`
	AS            string  `json:"as"`
	ASName        string  `json:"asname"`
	Reverse       string  `json:"reverse"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"`
	Hosting       bool    `json:"hosting"`
	Query         string  `json:"query"`
}

func (r *IPAPIClient) GetPublicIP(ctx context.Context) (*IPQueryResponse, error) {
	return r.queryIP(ctx, r.endpoint)
}

func (r *IPAPIClient) QueryIPAddress(ctx context.Context, ip string) (*IPQueryResponse, error) {
	return r.queryIP(ctx, r.endpoint+ip)
}

func (r *IPAPIClient) queryIP(ctx context.Context, endpoint string) (*IPQueryResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?fields=status,message,continent,continentCode,country,countryCode,region,regionName,city,district,zip,lat,lon,timezone,offset,currency,isp,org,as,asname,reverse,mobile,proxy,hosting,query", nil)
	if err != nil {
		return nil, err
	}

	response, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var ipQueryResponse IPQueryResponse
	err = json.NewDecoder(response.Body).Decode(&ipQueryResponse)
	if err != nil {
		return nil, err
	}

	return &ipQueryResponse, nil
}
