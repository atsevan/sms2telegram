package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sms2telegram/sms"
	"sms2telegram/telegram"
	"testing"
	"time"
)

type MockSmsReader struct {
	messages []sms.Sms
	err      error
}

func (m *MockSmsReader) ReadSMS() ([]sms.Sms, error) {
	return m.messages, m.err
}

type MockTelegramSender struct {
	err error
}

func (m *MockTelegramSender) SendMessage(message string) error {
	return m.err
}

func TestPollSMS(t *testing.T) {
	// Create a test server to mock the Telegram API endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"ok":true,"result":{"message_id":1}}`))
		// Verify the request parameters
		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()

	gc := &sms.GammuClient{
		HTTPClient: ts.Client(),
		Endpoint:   ts.URL,
		Username:   "admin",
		Password:   "password",
	}

	tc := &telegram.TelegramClient{
		HTTPClient: ts.Client(),
		ChatID:     "telegram-chat-id",
		Token:      "telegram-token",
		URL:        ts.URL + "/bottelegram-token/sendMessage",
	}

	// Start the polling in a separate goroutine
	go PollSMS(context.Background(), gc, tc, 50*time.Millisecond)

	// Wait for a few seconds to allow polling to occur
	time.Sleep(time.Second)
}

func TestPollSMSFetchAttempts(t *testing.T) {
	mockSmsReader := &MockSmsReader{
		err: errors.New("fetch error"),
	}

	mockTelegramSender := &MockTelegramSender{}

	t1 := time.Now()
	PollSMS(context.Background(), mockSmsReader, mockTelegramSender, 50*time.Millisecond)
	if time.Since(t1) > 300*time.Millisecond {
		t.Fatalf("Expected polling exists after 300ms, got %v", time.Since(t1))
	}
	if time.Since(t1) < 250*time.Millisecond {
		t.Fatalf("Expected polling takes more then 250ms, got %v", time.Since(t1))
	}
}

func TestPollSMSNoNewMessages(t *testing.T) {
	mockSmsReader := &MockSmsReader{
		err: sms.ErrNoNewMessages,
	}

	mockTelegramSender := &MockTelegramSender{}

	// Start the polling in a separate goroutine
	go PollSMS(context.Background(), mockSmsReader, mockTelegramSender, 50*time.Millisecond)

	// Wait for a few seconds to allow polling to occur
	time.Sleep(time.Second)
}

func TestPollSMSSuccess(t *testing.T) {
	mockSmsReader := &MockSmsReader{
		messages: []sms.Sms{
			{ID: "1", Number: "12345", Text: "Test message"},
		},
	}

	mockTelegramSender := &MockTelegramSender{}

	// Start the polling in a separate goroutine
	go PollSMS(context.Background(), mockSmsReader, mockTelegramSender, 50*time.Millisecond)

	// Wait for a few seconds to allow polling to occur
	time.Sleep(time.Second)

}
