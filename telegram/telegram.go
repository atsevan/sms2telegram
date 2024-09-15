package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TelegramClient represents a client to interact with the Telegram API.
type TelegramClient struct {
	HTTPClient *http.Client
	ChatID     string
	URL        string
	Token      string
}

// SendMessage sends a message to a Telegram chat.
func (t *TelegramClient) SendMessage(message string) error {
	payload := map[string]string{
		"chat_id": t.ChatID,
		"text":    message,
		// "parse_mode": parseMode,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, t.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s. %s", resp.Status, response["description"])
	}

	return nil
}
