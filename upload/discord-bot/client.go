package discord_bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"

	"github.com/yakiroren/dss-common/db"
	"github.com/yakiroren/dss-common/models"
)

type DiscordBotClient struct {
	dataStore     db.DataStore
	config        DiscordBotConfig
	session       *discordgo.Session
	totalChannels int
}

type DiscordBotConfig struct {
	DiscordStorageChannels []string `env:",required,notEmpty"`
	DiscordBotToken        string   `env:",required,notEmpty"`
}

func New(dataStore db.DataStore, config DiscordBotConfig) (*DiscordBotClient, error) {
	session, err := discordgo.New(config.DiscordBotToken)
	if err != nil {
		return nil, fmt.Errorf("faile to create discordgo client: %w", err)
	}

	session.UserAgent = ""
	session.MaxRestRetries = 5
	session.Client.Timeout = 3 * time.Minute

	return &DiscordBotClient{config: config, dataStore: dataStore, session: session, totalChannels: len(config.DiscordStorageChannels)}, nil
}

func (c *DiscordBotClient) Upload(ctx context.Context, id string, file []byte, fragmentID string) error {
	fragment, err := c.upload(fragmentID, bytes.NewBuffer(file))
	if err != nil {
		return err
	}

	return c.dataStore.AppendFragment(ctx, id, *fragment)
}

func (c *DiscordBotClient) upload(fragmentName string, file io.Reader) (*models.Fragment, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	channelID := c.config.DiscordStorageChannels[r1.Intn(c.totalChannels-1)]

	log.Debug("chosen channel id ", channelID)

	resp, err := c.session.ChannelFileSend(channelID, fragmentName, file)
	if err != nil {
		return nil, err
	}

	return &models.Fragment{
		ChannelID: resp.ChannelID,
		MessageID: resp.Attachments[0].ID,
		Name:      resp.Attachments[0].Filename,
		Size:      resp.Attachments[0].Size,
	}, nil
}
