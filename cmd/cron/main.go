package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/awlsring/dynamic-ip-watcher/internal/adapters/primary/watcher"
	cloudflare_dns_updater "github.com/awlsring/dynamic-ip-watcher/internal/adapters/secondary/dns_updater/cloudflare"
	ipapi_ip_retriever "github.com/awlsring/dynamic-ip-watcher/internal/adapters/secondary/ip_retriever/ip_api"
	"github.com/awlsring/dynamic-ip-watcher/internal/adapters/secondary/notifier/discord_webhook"
	local_storage "github.com/awlsring/dynamic-ip-watcher/internal/adapters/secondary/storage/local"
	"github.com/awlsring/dynamic-ip-watcher/internal/config"
	"github.com/awlsring/dynamic-ip-watcher/internal/core/service/address"
	ipapi "github.com/awlsring/dynamic-ip-watcher/internal/pkg/ip-api"
	"github.com/awlsring/dynamic-ip-watcher/internal/ports/gateway"
	"github.com/cloudflare/cloudflare-go"
	"github.com/rs/zerolog/log"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadNotifiers(cfg *config.Config) []gateway.Notifier {
	var notifiers []gateway.Notifier
	for _, notifier := range cfg.Notifiers {
		switch notifier {
		case notifier.(config.DiscordNotifierConfig):
			notifierConfig := notifier.(config.DiscordNotifierConfig)
			notifier := discord_webhook.New(
				notifierConfig.WebhookUrl,
				notifierConfig.AvatarUrl,
				notifierConfig.Username,
				http.DefaultClient,
			)
			notifiers = append(notifiers, notifier)
		default:
			log.Warn().Msgf("Unknown notifier type: %s", notifier.GetNotifierType())
		}
	}
	return notifiers
}

func loadDnsUpdater(cfg *config.Config) gateway.DNSUpdater {
	switch cfg.DNSRecord.Type {
	case config.DnsRecordTypeCloudflare:
		cloudflareClient, err := cloudflare.NewWithAPIToken(cfg.DNSRecord.APIKey)
		panicOnError(err)
		return cloudflare_dns_updater.New(cfg.DNSRecord.ZoneName, cfg.DNSRecord.RecordName, cloudflareClient)
	default:
		log.Warn().Msgf("Unknown DNS updater type: %s", cfg.DNSRecord.Type)
		return nil
	}
}

func loadIpRetriever(*config.Config) gateway.IPRetriever {
	ipApiClient := ipapi.New()
	ipRetriever := ipapi_ip_retriever.New(ipApiClient)
	return ipRetriever
}

func loadStorage(cfg *config.Config) gateway.Storage {
	return local_storage.New(cfg.Storage.Directory)
}

func main() {
	cfg, err := config.Load()
	panicOnError(err)

	notifiers := loadNotifiers(cfg)
	dnsUpdater := loadDnsUpdater(cfg)
	ipRetriever := loadIpRetriever(cfg)
	storage := loadStorage(cfg)

	addressService := address.NewService(dnsUpdater, ipRetriever, notifiers, storage)

	watcher := watcher.New(addressService)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error)
	go func() {
		err := watcher.Run(ctx)
		if err != nil {
			done <- err
			return
		}

		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error().Err(err).Msg("Failed to complete")
			os.Exit(1)
		}
		log.Info().Msg("Completed successfully")
	case <-ctx.Done():
		log.Warn().Msg("Timeout reached before completion, exiting...")
		os.Exit(1)
	case <-sigs:
		log.Warn().Msg("Received signal, exiting...")
		os.Exit(1)
	}
}
