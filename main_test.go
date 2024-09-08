package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetchSMS(t *testing.T) {
	sms := Sms{
		Date:   "2022-01-01",
		Number: "+123456789",
		State:  "received",
		Text:   "Test SMS",
	}

	// Create a test server to mock the sms-gammu-gateway endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request parameters
		if r.Method != http.MethodGet {
			t.Errorf("fetchSMS received incorrect HTTP method: got %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/getsms" {
			t.Errorf("fetchSMS received incorrect URL path: got %s, want %s", r.URL.Path, "/getsms")
		}

		// Verify the Authorization header
		// admin:password base64 encoded is YWRtaW46cGFzc3dvcmQ=
		if r.Header.Get("Authorization") != "Basic YWRtaW46cGFzc3dvcmQ=" {
			t.Errorf("fetchSMS received incorrect Authorization header: got %s, want %s", r.Header.Get("Authorization"), "Basic YWRtaW46cGFzc3dvcmQ=")
		}
		// Return a sample SMS JSON response
		sms := sms
		json.NewEncoder(w).Encode(sms)
	}))
	defer ts.Close()

	// Create a new Sms2Telegram instance with the test server URL
	s := &GammuClient{
		HTTPClient: ts.Client(),
		Endpoint:   ts.URL,
		Username:   "admin",
		Password:   "password",
	}

	// Call the fetchSMS method
	smsTest, err := s.FetchSMS()

	// Verify the result
	if err != nil {
		t.Errorf("fetchSMS returned an error: %v", err)
	}
	if sms.Date != "2022-01-01" {
		t.Errorf("fetchSMS returned incorrect date: got %s, want %s", smsTest.Date, sms.Date)
	}
	if sms.Number != "+123456789" {
		t.Errorf("fetchSMS returned incorrect number: got %s, want %s", smsTest.Number, sms.Number)
	}
	if sms.State != "received" {
		t.Errorf("fetchSMS returned incorrect state: got %s, want %s", smsTest.State, sms.State)
	}
	if sms.Text != "Test SMS" {
		t.Errorf("fetchSMS returned incorrect text: got %s, want %s", smsTest.Text, sms.Text)
	}
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
	err := tc.sendTelegramMessage("Test message")

	// Verify the result
	if err != nil {
		t.Errorf("sendTelegramMessage returned an error: %v", err)
	}
}

func TestPollSMS(t *testing.T) {
	// Create a stop channel to stop the polling
	stopChan := make(chan struct{})
	defer close(stopChan)

	gc := &GammuClient{
		HTTPClient: http.DefaultClient,
		Endpoint:   "http://localhost:8080/",
		Username:   "admin",
		Password:   "password",
	}

	tc := &TelegramClient{
		HTTPClient: http.DefaultClient,
		ChatID:     "telegram-chat-id",
		Token:      "telegram-token",
		URL:        "http://localhost:8080/bottelegram-token/sendMessage",
	}

	// Start the polling in a separate goroutine
	go PollSMS(*gc, *tc, time.Second, stopChan)

	// Wait for a few seconds to allow polling to occur
	time.Sleep(2 * time.Second)

	// Stop the polling by sending a signal to the stop channel
	stopChan <- struct{}{}
}

func TestReset(t *testing.T) {
	// Create a test server to mock the sms-gammu-gateway reset endpoint
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request parameters
		if r.Method != http.MethodGet {
			t.Errorf("Reset received incorrect HTTP method: got %s, want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/reset" {
			t.Errorf("Reset received incorrect URL path: got %s, want %s", r.URL.Path, "/reset")
		}

		// Verify the Authorization header
		// admin:password base64 encoded is YWRtaW46cGFzc3dvcmQ=
		if r.Header.Get("Authorization") != "Basic YWRtaW46cGFzc3dvcmQ=" {
			t.Errorf("Reset received incorrect Authorization header: got %s, want %s", r.Header.Get("Authorization"), "Basic YWRtaW46cGFzc3dvcmQ=")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create a new GammuClient instance with the test server URL
	gc := &GammuClient{
		HTTPClient: ts.Client(),
		Endpoint:   ts.URL,
		Username:   "admin",
		Password:   "password",
	}

	// Call the Reset method
	err := gc.Reset()

	// Verify the result
	if err != nil {
		t.Errorf("Reset returned an error: %v", err)
	}
}
