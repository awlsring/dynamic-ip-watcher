package cloudflare_dns_updater

import (
	"context"
	"net"

	"github.com/awlsring/dynamic-ip-watcher/internal/pkg/interfaces"
	"github.com/awlsring/dynamic-ip-watcher/internal/ports/gateway"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

const (
	A             = "A"
	RecordComment = "dynamic-ip-watcher"
)

type CloudflareDNSUpdater struct {
	zoneName string
	zoneId   string
	dnsName  string
	client   interfaces.CloudflareAPI
}

func New(zoneName, dnsName string, client interfaces.CloudflareAPI) gateway.DNSUpdater {
	return &CloudflareDNSUpdater{
		zoneName: zoneName,
		dnsName:  dnsName,
		client:   client,
	}
}

func (a *CloudflareDNSUpdater) RecordName() string {
	return a.dnsName
}

func (a *CloudflareDNSUpdater) GetRecordIpAddress(ctx context.Context) (net.IP, error) {
	record, err := a.describeRecord(ctx)
	if err != nil {
		return nil, err
	}

	return net.ParseIP(record.Content), nil
}

func (a *CloudflareDNSUpdater) CreateRecordWithIpAddress(ctx context.Context, ip net.IP) error {
	zoneId, err := a.getZoneId(ctx)
	if err != nil {
		return err
	}

	_, err = a.client.CreateDNSRecord(context.Background(), &cloudflare.ResourceContainer{Identifier: zoneId}, cloudflare.CreateDNSRecordParams{
		Type:    A,
		Name:    a.dnsName,
		Content: ip.String(),
		Comment: RecordComment,
	})

	return err
}

func (a *CloudflareDNSUpdater) UpdateRecordIpAddress(ctx context.Context, ip net.IP) error {
	zoneId, err := a.getZoneId(ctx)
	if err != nil {
		return err
	}

	cloudflareRecord, err := a.describeRecord(ctx)
	if err != nil {
		return err
	}

	_, err = a.client.UpdateDNSRecord(context.Background(), &cloudflare.ResourceContainer{Identifier: zoneId}, cloudflare.UpdateDNSRecordParams{
		ID:      cloudflareRecord.ID,
		Content: ip.String(),
	})

	return err
}

func (a *CloudflareDNSUpdater) getZoneId(context.Context) (string, error) {
	if a.zoneId == "" {
		zoneId, err := a.client.ZoneIDByName(a.zoneName)
		if err != nil {
			log.Error().Str("ZoneName", a.zoneName).Err(err).Msg("Failed to get zone ID by name")
			return "", err
		}
		a.zoneId = zoneId
	}

	return a.zoneId, nil
}

func (a *CloudflareDNSUpdater) describeRecord(ctx context.Context) (cloudflare.DNSRecord, error) {
	zoneId, err := a.client.ZoneIDByName(a.zoneName)
	if err != nil {
		return cloudflare.DNSRecord{}, err
	}

	records, _, err := a.client.ListDNSRecords(ctx, &cloudflare.ResourceContainer{Identifier: zoneId}, cloudflare.ListDNSRecordsParams{
		Type: A,
		Name: a.dnsName,
	})
	if err != nil {
		return cloudflare.DNSRecord{}, err
	}

	if len(records) == 0 {
		return cloudflare.DNSRecord{}, gateway.ErrRecordNotFound
	}

	if len(records) > 1 {
		return cloudflare.DNSRecord{}, gateway.ErrMultipleRecordsFound
	}

	return records[0], nil
}
