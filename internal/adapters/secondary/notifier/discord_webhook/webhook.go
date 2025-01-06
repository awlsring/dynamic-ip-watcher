package discord_webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/awlsring/dynamic-ip-watcher/internal/core/domain/event"
	"github.com/rs/zerolog/log"
)

const (
	Username = "DynamicIPWatcher"
)

type DiscordWebhookNotifier struct {
	client    *http.Client
	url       string
	avatarUrl string
	username  string
}

func New(url, avatarUrl, username string, client *http.Client) *DiscordWebhookNotifier {
	notifier := &DiscordWebhookNotifier{
		client:    client,
		url:       url,
		avatarUrl: avatarUrl,
		username:  username,
	}

	if notifier.username == "" {
		notifier.username = Username
	}

	return notifier
}

func (d *DiscordWebhookNotifier) SendEventMessage(ctx context.Context, event event.Event) error {
	if event == nil {
		log.Warn().Msg("event provided was empty, not sending.")
		return nil
	}

	discordMessage := DiscordWebhookMessage{
		Username: Username,
		Content:  event.AsMessage(),
	}

	payloadBytes, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
