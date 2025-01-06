package discord_webhook

type DiscordEmbed struct {
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Color       int    `json:"color,omitempty"`
	Footer      struct {
		Text    string `json:"text,omitempty"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"footer,omitempty"`
	Image struct {
		URL string `json:"url,omitempty"`
	} `json:"image,omitempty"`
	Thumbnail struct {
		URL string `json:"url,omitempty"`
	} `json:"thumbnail,omitempty"`
	Video struct {
		URL string `json:"url,omitempty"`
	} `json:"video,omitempty"`
	Provider struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"provider,omitempty"`
	Author struct {
		Name    string `json:"name,omitempty"`
		URL     string `json:"url,omitempty"`
		IconURL string `json:"icon_url,omitempty"`
	} `json:"author,omitempty"`
	Fields []struct {
		Name   string `json:"name,omitempty"`
		Value  string `json:"value,omitempty"`
		Inline bool   `json:"inline,omitempty"`
	} `json:"fields,omitempty"`
}

// https://discord.com/developers/docs/resources/webhook#execute-webhook
type DiscordWebhookMessage struct {
	Content   string         `json:"content,omitempty"`
	Username  string         `json:"username,omitempty"`
	AvatarUrl string         `json:"avatar_url,omitempty"`
	TTS       bool           `json:"tts,omitempty"`
	Embeds    []DiscordEmbed `json:"embeds,omitempty"`
}
