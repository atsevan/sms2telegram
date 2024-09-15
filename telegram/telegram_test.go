package telegram

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func emptyJson() []byte {
	// Create an empty JSON object
	emptyJSON := struct{}{}

	// Marshal the empty JSON object to bytes
	bytes, _ := json.Marshal(emptyJSON)
	return bytes
}

func TestSendTelegramMessage(t *testing.T) {

	// Create a test server to mock the Telegram API endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request parameters
		if r.Method != http.MethodPost {
			t.Errorf("sendTelegramMessage received incorrect HTTP method: got %s, want %s", r.Method, http.MethodPost)
		}
		if r.URL.Path != "/bottelegram-token/sendMessage" {
			t.Errorf("sendTelegramMessage received incorrect URL path: got %s, want %s", r.URL.Path, "/bottelegram-token/sendMessage")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("sendTelegramMessage should set application/json as a content-type")
		}

		w.Write(emptyJson())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new TelegramClient instance with the test server URL
	tc := &TelegramClient{
		HTTPClient: ts.Client(),
		ChatID:     "telegram-chat-id",
		Token:      "telegram-token",
		URL:        ts.URL + "/bottelegram-token/sendMessage",
	}

	// Call the sendTelegramMessage method
	err := tc.SendMessage("Test message")

	// Verify the result
	if err != nil {
		t.Errorf("sendTelegramMessage returned an error: %v", err)
	}
}

func TestSendTelegramMessageError(t *testing.T) {

	// Create a test server to mock the Telegram API endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"description": "Internal Server Error"}`))
	}))
	defer ts.Close()

	// Create a new TelegramClient instance with the test server URL
	tc := &TelegramClient{
		HTTPClient: ts.Client(),
		ChatID:     "telegram-chat-id",
		Token:      "telegram-token",
		URL:        ts.URL + "/bottelegram-token/sendMessage",
	}

	// Call the sendTelegramMessage method
	err := tc.SendMessage("Test message")

	// Verify the result
	if err == nil {
		t.Errorf("sendTelegramMessage should have returned an error")
	}
}

func TestSendTelegramMessageInvalidJSON(t *testing.T) {

	// Create a test server to mock the Telegram API endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer ts.Close()

	// Create a new TelegramClient instance with the test server URL
	tc := &TelegramClient{
		HTTPClient: ts.Client(),
		ChatID:     "telegram-chat-id",
		Token:      "telegram-token",
		URL:        ts.URL + "/bottelegram-token/sendMessage",
	}

	// Call the sendTelegramMessage method
	err := tc.SendMessage("Test message")

	// Verify the result
	if err == nil {
		t.Errorf("sendTelegramMessage should have returned an error due to invalid JSON response")
	}
}
