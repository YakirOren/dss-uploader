package webhooks

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"net/textproto"

	log "github.com/sirupsen/logrus"
	"github.com/yakiroren/dss-common/models"

	"github.com/yakiroren/dss-common/db"
)

type DiscordWebhookClient struct {
	webhookURLs []string
	dataStore   db.DataStore
}

type DiscordWebhookConfig struct {
	DiscordWebhooks []string `env:",required,notEmpty"`
}

func New(dataStore db.DataStore, config DiscordWebhookConfig) (*DiscordWebhookClient, error) {
	return &DiscordWebhookClient{webhookURLs: config.DiscordWebhooks, dataStore: dataStore}, nil
}

func (c *DiscordWebhookClient) Upload(ctx context.Context, path string, file []byte, fragmentID string) error {
	index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(c.webhookURLs))))

	webhookURL := c.webhookURLs[index.Int64()]

	resp, err := c.upload(ctx, fragmentID, file, webhookURL)
	if err != nil {
		return err
	}

	fragment := models.Fragment{
		ChannelID: resp.ChannelID,
		MessageID: resp.Attachments[0].ID,
		Name:      resp.Attachments[0].Filename,
		Size:      resp.Attachments[0].Size,
	}

	err = c.dataStore.AppendFragment(ctx, path, fragment)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// MultipartBodyWithJSON returns the contentType and body for a discord request
// data  : The object to encode for payload_json in the multipart request
// files : Files to include in the request.
func MultipartBodyWithJSON(data interface{}, filename string, file []byte) (string, []byte, error) {
	body := &bytes.Buffer{}
	bodywriter := multipart.NewWriter(body)

	payload, err := json.Marshal(data)
	if err != nil {
		return "", nil, err
	}

	var p io.Writer

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", "application/json")

	p, err = bodywriter.CreatePart(h)
	if err != nil {
		return "", nil, err
	}

	if _, err = p.Write(payload); err != nil {
		return "", nil, err
	}

	h2 := make(textproto.MIMEHeader)
	h2.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, 1, filename))

	h2.Set("Content-Type", "application/octet-stream")

	p, err = bodywriter.CreatePart(h2)
	if err != nil {
		return "", nil, err
	}

	_, _ = p.Write(file)

	err = bodywriter.Close()
	if err != nil {
		return "", nil, err
	}

	return bodywriter.FormDataContentType(), body.Bytes(), nil
}

func (c *DiscordWebhookClient) upload(ctx context.Context,
	filename string,
	file []byte, webhookURL string,
) (*DiscordResponse, error) {
	contentType, body, err := MultipartBodyWithJSON(struct{}{}, filename, file)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP POST request with the webhook URL.
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("Content-Type", contentType)

	// Send the HTTP request.
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer res.Body.Close()

	// Check if the response status code is not successful.
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", res.StatusCode)
	}

	result := &DiscordResponse{}
	if err = json.NewDecoder(res.Body).Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
