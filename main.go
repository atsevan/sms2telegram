package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	httpTimeout     = 10 * time.Second
	sendMessageTmpl = "https://api.telegram.org/bot%s/sendMessage"
)

// errNoNewMessages represents the error when there are no new messages.
var errNoNewMessages = fmt.Errorf("no new messages")

var (
	endpoint       = flag.String("endpoint", getEnv("ENDPOINT", "http://localhost:5000/getsms"), "sms-gammu-gateway URL")
	username       = flag.String("username", getEnv("USERNAME", "admin"), "sms-gammu-gateway username")
	password       = flag.String("password", getEnv("PASSWORD", "password"), "sms-gammu-gateway password")
	telegramToken  = flag.String("telegram-token", getEnv("TELEGRAM_TOKEN", ""), "Telegram bot token")
	telegramChatID = flag.String("telegram-chat-id", getEnv("TELEGRAM_CHAT_ID", ""), "Telegram chat ID")
	interval       = flag.Duration("interval", getEnvDuration("INTERVAL", 5*time.Second), "polling interval")
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// Sms represents a SMS message.
type Sms struct {
	Date   string `json:"Date"`
	Number string `json:"number"`
	State  string `json:"State"`
	Text   string `json:"Text"`
}

func (s Sms) String() string {
	return fmt.Sprintf("%s sent on %s (%s)\n%s", s.Number, s.Date, s.State, s.Text)
}

// Validate validates the Sms struct.
func (s Sms) Validate() error {
	if s.Date == "" && s.Number == "" && s.State == "" && s.Text == "" {
		return errNoNewMessages
	}

	if s.Date == "" {
		return fmt.Errorf("Date is required")
	}
	if s.Number == "" {
		return fmt.Errorf("Number is required")
	}
	if s.State == "" {
		return fmt.Errorf("State is required")
	}
	if s.Text == "" {
		return fmt.Errorf("Text is required")
	}
	return nil
}

// GammuClient represents a client to interact with the sms-gammu-gateway.
type GammuClient struct {
	HTTPClient *http.Client
	GetSMSUrl  string
	Username   string
	Password   string
}

// TelegramClient represents a client to interact with the Telegram API.
type TelegramClient struct {
	HTTPClient *http.Client
	ChatID     string
	URL        string
	Token      string
}

// fetchSMS fetches the SMS from the sms-gammu-gateway endpoint.
func (g *GammuClient) fetchSMS() (Sms, error) {
	req, err := http.NewRequest(http.MethodGet, g.GetSMSUrl, nil)
	if err != nil {
		return Sms{}, err
	}

	req.SetBasicAuth(g.Username, g.Password)

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return Sms{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Sms{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sms{}, err
	}

	var sms Sms
	err = json.Unmarshal(body, &sms)
	if err != nil {
		return Sms{}, err
	}

	if err = sms.Validate(); err != nil {
		return Sms{}, err
	}

	return sms, nil
}

// sendTelegramMessage sends a message to a Telegram chat.
func (t *TelegramClient) sendTelegramMessage(message string) error {
	payload := map[string]string{
		"chat_id": t.ChatID,
		"text":    message,
		// "parse_mode": "Markdown",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, t.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}

// pollSMS polls the SMS endpoint and sends a Telegram message for each new SMS.
func PollSMS(gummuc GammuClient, telegramc TelegramClient, interval time.Duration, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping the polling.")
			return
		default:
			// Fetch SMS from endpoint
			sms, err := gummuc.fetchSMS()
			if err != nil {
				if err != errNoNewMessages {
					log.Println("Error fetching SMS:", err)
				}
				time.Sleep(interval)
				continue
			}

			log.Println("Got new sms: ", sms)

			// Send Telegram message
			err = telegramc.sendTelegramMessage(sms.String())
			if err != nil {
				log.Println("Error sending Telegram message:", err)
			}

			// Wait for the next interval
			time.Sleep(interval)
		}
	}
}

func main() {
	flag.Parse()

	if *telegramToken == "" || *telegramChatID == "" {
		log.Println("Telegram token and chat ID are required")
		os.Exit(1)
	}

	// Create clients
	gummuc := GammuClient{
		HTTPClient: &http.Client{Timeout: httpTimeout},
		GetSMSUrl:  *endpoint,
		Username:   *username,
		Password:   *password,
	}
	telegramc := TelegramClient{
		HTTPClient: &http.Client{Timeout: httpTimeout},
		ChatID:     *telegramChatID,
		URL:        fmt.Sprintf(sendMessageTmpl, *telegramToken),
		Token:      *telegramToken,
	}

	// Create a channel to stop polling
	stopChan := make(chan struct{})

	// Start polling in a separate goroutine
	go PollSMS(gummuc, telegramc, *interval, stopChan)

	// Wait for a signal to stop polling
	select {}

	// Stop polling
	// close(stopChan)
}
