package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
)

const (
	DefaultConfigPath = "/var/lib/dynamic-ip-watcher/config.json"
	DefaultStorageDir = "/var/lib/dynamic-ip-watcher"

	ConfigPathEnvVar     = "CONFIG_PATH"
	ZoneIDEnvVar         = "ZONE_ID"
	RecordNameEnvVar     = "RECORD_NAME"
	DiscordWebhookEnvVar = "DISCORD_WEBHOOK"
	LocalStorageDirEnv   = "LOCAL_STORAGE_DIR"
)

const (
	NotifierTypeDiscord = "discord"
)

const (
	DnsRecordTypeCloudflare = "cloudflare"
)

type Notifier interface {
	GetNotifierType() string
}

type DiscordNotifierConfig struct {
	Type       string `json:"type"`
	WebhookUrl string `json:"webhookUrl"`
	Username   string `json:"username"`
	AvatarUrl  string `json:"avatarUrl"`
}

func (d DiscordNotifierConfig) GetNotifierType() string {
	return d.Type
}

type NotifierConfig struct {
	Type     string `json:"type"`
	Endpoint string `json:"endpoint"`
}

type DNSRecordConfig struct {
	Type       string `json:"type"`
	APIKey     string `json:"apiKey"`
	ZoneName   string `json:"zoneName"`
	RecordName string `json:"recordName"`
}

type StorageConfig struct {
	Directory string `json:"directory"`
}

type Config struct {
	DNSRecord DNSRecordConfig `json:"dnsRecord"`
	Storage   StorageConfig   `json:"storage"`
	Notifiers []Notifier      `json:"notifiers"`
}

func getEnvOrDefault(envVar, defaultValue string) string {
	value := os.Getenv(envVar)
	if value == "" {
		return defaultValue
	}

	return value
}

// check if a --config-path arg was passed in,
// if not, check if the CONFIG_PATH env var is set
// if not, use the default path
func determineConfigPath() string {
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--config-path" && i+1 < len(args) {
			return args[i+1]
		}
	}

	value := os.Getenv(ConfigPathEnvVar)
	if value == "" {
		return DefaultConfigPath
	}

	return value
}

func loadConfigFromFile(cfg *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	var rawConfig struct {
		DNSRecord DNSRecordConfig   `json:"dnsRecord"`
		Storage   StorageConfig     `json:"storage"`
		Notifiers []json.RawMessage `json:"notifiers"`
	}

	if err := json.NewDecoder(file).Decode(&rawConfig); err != nil {
		return err
	}

	if rawConfig.Storage.Directory == "" {
		rawConfig.Storage.Directory = DefaultStorageDir
	}

	cfg.DNSRecord = rawConfig.DNSRecord
	cfg.Storage = rawConfig.Storage

	for _, rawNotifier := range rawConfig.Notifiers {
		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawNotifier, &base); err != nil {
			return err
		}

		var notifier Notifier
		switch base.Type {
		case NotifierTypeDiscord:
			var discordConfig DiscordNotifierConfig
			if err := json.Unmarshal(rawNotifier, &discordConfig); err != nil {
				return err
			}
			replaceFilePaths(&discordConfig)
			notifier = discordConfig
		default:
			return errors.New("unknown notifier type: " + base.Type)
		}

		cfg.Notifiers = append(cfg.Notifiers, notifier)
	}

	return replaceFilePaths(cfg)
}

func replaceFilePaths(structs ...interface{}) error {
	for _, s := range structs {
		v := reflect.ValueOf(s)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)

			if field.Kind() == reflect.String && field.String() != "" {
				if isFilePath(field.String()) {
					content, err := tryReadFile(field.String())
					if err == nil {
						field.SetString(content)
					}
				}
			}

			if field.Kind() == reflect.Struct {
				if err := replaceFilePaths(field.Addr().Interface()); err != nil {
					return err
				}
			}

			if field.Kind() == reflect.Slice && field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					elem := field.Index(j).Addr().Interface()
					if err := replaceFilePaths(elem); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func isFilePath(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func tryReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func setConfigFromEnv(cfg *Config) {
	cfg.DNSRecord.ZoneName = getEnvOrDefault(ZoneIDEnvVar, cfg.DNSRecord.ZoneName)
	cfg.DNSRecord.RecordName = getEnvOrDefault(RecordNameEnvVar, cfg.DNSRecord.RecordName)

	cfg.Storage.Directory = getEnvOrDefault(LocalStorageDirEnv, cfg.Storage.Directory)

	if os.Getenv(DiscordWebhookEnvVar) != "" {
		if len(cfg.Notifiers) == 0 {
			cfg.Notifiers = make([]Notifier, 0)
		}
		didSet := setDiscordWebhookValue(cfg, os.Getenv(DiscordWebhookEnvVar))
		if !didSet {
			cfg.Notifiers = append(cfg.Notifiers, DiscordNotifierConfig{
				Type:       NotifierTypeDiscord,
				WebhookUrl: os.Getenv(DiscordWebhookEnvVar),
			})
		}
	}
}

func setDiscordWebhookValue(cfg *Config, value string) bool {
	for i, notifier := range cfg.Notifiers {
		if notifier.GetNotifierType() == NotifierTypeDiscord {
			discordNotifier := notifier.(DiscordNotifierConfig)
			discordNotifier.WebhookUrl = value
			cfg.Notifiers[i] = discordNotifier
			return true
		}
	}

	return false
}

// LoadConfig attempts to load the configuration from a path, then sets values from Environment Variables.
func Load() (*Config, error) {
	cfg := &Config{}

	configPath := determineConfigPath()
	log.Info().Msgf("Loading configuration from %s", configPath)

	err := loadConfigFromFile(cfg, configPath)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load configuration")
		return nil, err
	}

	setConfigFromEnv(cfg)

	return cfg, nil
}
