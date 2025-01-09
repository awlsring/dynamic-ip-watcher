package address

import (
	"context"
	"fmt"

	"github.com/awlsring/dynamic-ip-watcher/internal/core/domain/event"
	"github.com/awlsring/dynamic-ip-watcher/internal/ports/gateway"
	"github.com/awlsring/dynamic-ip-watcher/internal/ports/service"
	"github.com/rs/zerolog/log"
)

type Service struct {
	dnsUpdater  gateway.DNSUpdater
	ipRetriever gateway.IPRetriever
	notifiers   []gateway.Notifier
	storage     gateway.Storage
}

func NewService(dnsUpdater gateway.DNSUpdater, ipRetriever gateway.IPRetriever, notifiers []gateway.Notifier, storage gateway.Storage) service.Address {
	return &Service{
		dnsUpdater:  dnsUpdater,
		ipRetriever: ipRetriever,
		notifiers:   notifiers,
		storage:     storage,
	}
}

func (s *Service) sendEventToNotifiers(ctx context.Context, event event.Event) {
	log.Info().Msg("Sending event to notifiers")
	for _, notifier := range s.notifiers {
		err := notifier.SendEventMessage(ctx, event)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to send event message")
		}
	}
}

func (s *Service) DetectAndHandleAddressChange(ctx context.Context) error {
	var eventMessage event.Event
	defer func() {
		s.sendEventToNotifiers(ctx, eventMessage)
	}()

	log.Info().Msg("Detecting IP address change")
	previousIP, err := s.storage.GetLastKnownIPAddress(ctx)
	if err != nil {
		eventMessage = event.NewFailedUpdateEvent("Failed to determine the last known IP address", err)
		log.Error().Err(err).Msg("Failed to get last known IP address")
		return err
	}
	log.Info().Str("previous_ip", previousIP.String()).Msg("Previous IP address")

	log.Info().Msg("Retrieving current IP address")
	currentIP, err := s.ipRetriever.GetPublicIPv4(ctx)
	if err != nil {
		eventMessage = event.NewFailedUpdateEvent("Failed to determine current IP address", err)
		log.Error().Err(err).Msg("Failed to get current IP address")
		return err
	}
	log.Info().Str("current_ip", currentIP.String()).Msg("Current IP address")

	if previousIP.Equal(currentIP) {
		log.Info().Msg("IP address has not changed")
		return nil
	}
	log.Info().Msg("IP address has changed")
	eventMessage = event.NewChangeEvent("IP address changed from " + previousIP.String() + " to " + currentIP.String())

	log.Info().Msg("Saving current IP address")
	err = s.storage.SaveIPAddress(ctx, currentIP)
	if err != nil {
		eventMessage = event.NewFailedUpdateEvent("Failed to store new IP address", err)
		log.Error().Err(err).Msg("Failed to save current IP address")
		return err
	}

	if s.dnsUpdater != nil {
		log.Info().Msg("Updating DNS A record with new IP address")
		err = s.dnsUpdater.UpdateRecordIpAddress(ctx, currentIP)
		if err != nil {
			eventMessage = event.NewFailedUpdateEvent("Failed to update DNS Record with new IP address", err)
			log.Error().Err(err).Msg("Failed to update DNS A record")
			return err
		}
		message := fmt.Sprintf("IP address changed from %s to %s. DNS Record %s updated with new address.", previousIP.String(), currentIP.String(), s.dnsUpdater.RecordName())
		eventMessage = event.NewChangeEvent(message)
	}

	return nil
}
