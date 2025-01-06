package interfaces

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

// implements cloudflare.API
type CloudflareAPI interface {
	CreateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.CreateDNSRecordParams) (cloudflare.DNSRecord, error)
	GetDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, recordID string) (cloudflare.DNSRecord, error)
	ListDNSRecords(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.ListDNSRecordsParams) ([]cloudflare.DNSRecord, *cloudflare.ResultInfo, error)
	UpdateDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, params cloudflare.UpdateDNSRecordParams) (cloudflare.DNSRecord, error)
	DeleteDNSRecord(ctx context.Context, rc *cloudflare.ResourceContainer, recordID string) error
	ZoneIDByName(zoneName string) (string, error)
}
