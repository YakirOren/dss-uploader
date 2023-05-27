package discord_webhooks

import "time"

type DiscordResponse struct {
	Id        string `json:"id"`
	Type      int    `json:"type"`
	Content   string `json:"content"`
	ChannelId string `json:"channel_id"`
	Author    struct {
		Bot           bool        `json:"bot"`
		Id            string      `json:"id"`
		Username      string      `json:"username"`
		Avatar        interface{} `json:"avatar"`
		Discriminator string      `json:"discriminator"`
	} `json:"author"`
	Attachments []struct {
		Id          string `json:"id"`
		Filename    string `json:"filename"`
		Size        int    `json:"size"`
		Url         string `json:"url"`
		ProxyUrl    string `json:"proxy_url"`
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		ContentType string `json:"content_type"`
	} `json:"attachments"`
	Embeds          []interface{} `json:"embeds"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	Pinned          bool          `json:"pinned"`
	MentionEveryone bool          `json:"mention_everyone"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Flags           int           `json:"flags"`
	Components      []interface{} `json:"components"`
	WebhookId       string        `json:"webhook_id"`
}
