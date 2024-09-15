package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sms2telegram/sms"
	"sms2telegram/telegram"
	"sync"
	"time"
)

const (
	httpTimeout     = 10 * time.Second
	sendMessageTmpl = "https://api.telegram.org/bot%s/sendMessage"
)

var (
	// GammyClient configuration
	endpoint = flag.String("endpoint", getEnv("ENDPOINT", "http://localhost:5000/"), "sms-gammu-gateway URL")
	username = flag.String("username", getEnv("USERNAME", "admin"), "sms-gammu-gateway username")
	password = flag.String("password", getEnv("PASSWORD", "password"), "sms-gammu-gateway password")

	// TelegramClient configuration
	telegramToken  = flag.String("telegram-token", getEnv("TELEGRAM_TOKEN", ""), "Telegram bot token")
	telegramChatID = flag.String("telegram-chat-id", getEnv("TELEGRAM_CHAT_ID", ""), "Telegram chat ID")

	// Polling configuration
	interval = flag.Duration("interval", getEnvDuration("INTERVAL", 5*time.Second), "polling interval")
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

type SmsReader interface {
	ReadSMS() ([]sms.Sms, error)
}

type TelegramSender interface {
	SendMessage(message string) error
}

// PollSMS polls the SMS endpoint and sends a Telegram message for each new SMS.
func PollSMS(ctx context.Context, sr SmsReader, ts TelegramSender, interval time.Duration) {
	fetchAttempts := 0
	ctx, cancel := context.WithCancel(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Context is done. Stopping the polling.")
			cancel()
			return
		default:
			// Fetch SMSs
			msgs, err := sr.ReadSMS()
			if err != nil {
				if err != sms.ErrNoNewMessages {
					log.Println("Error fetching SMS:", err)
					fetchAttempts++
				}
				if fetchAttempts > 5 {
					log.Println("Too many fetch attempts. Stopping the polling.")
					cancel()
					return
				}
				time.Sleep(interval)
				continue
			}

			log.Println("Got new sms: ", msgs)

			// Send a Telegram message for each SMS
			for _, s := range msgs {
				log.Println("SMS: ", s)
				err = ts.SendMessage(s.String())
				if err != nil {
					log.Println("Error sending Telegram message:", err)
				}
			}
			fetchAttempts = 0
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
	gc := sms.GammuClient{
		HTTPClient: &http.Client{Timeout: httpTimeout},
		Endpoint:   *endpoint,
		Username:   *username,
		Password:   *password,
	}

	tc := telegram.TelegramClient{
		HTTPClient: &http.Client{Timeout: httpTimeout},
		ChatID:     *telegramChatID,
		URL:        fmt.Sprintf(sendMessageTmpl, *telegramToken),
		Token:      *telegramToken,
	}

	// Send Telegram message
	err := tc.SendMessage("Bot has started")
	if err != nil {
		log.Fatalf("Error sending Telegram message: %v", err)
	}

	// Create a WaitGroup to wait for the goroutine to finish
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Start polling in a separate goroutine
		PollSMS(context.Background(), &gc, &tc, *interval)
	}()

	// Wait for the goroutine to finish
	wg.Wait()

}
