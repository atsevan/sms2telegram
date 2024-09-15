package sms

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGammuClient_ReadSMS(t *testing.T) {

	tests := []struct {
		name     string
		response string
		want     Sms
		wantErr  error
	}{
		{
			name:     "Test message with all fields",
			response: `{"ID":"1","Text":"Test message 1","Number":"+123456789","State":"received","Date":"2022-01-01"}`,
			want:     Sms{ID: "1", Text: "Test message 1", Number: "+123456789", State: "received", Date: "2022-01-01"},
			wantErr:  nil,
		},
		{
			name:     "Test message with missing date",
			response: `{"ID":"2","Text":"Test message 2","Number":"+123456789","State":"received"}`,
			want:     Sms{ID: "2", Text: "Test message 2", Number: "+123456789", State: "received"},
			wantErr:  ErrNoDate,
		},
		{
			name:     "Empty message",
			response: `{}`,
			want:     Sms{},
			wantErr:  ErrNoNewMessages,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/getsms" {
					t.Fatalf("Expected to request '/getsms', got: %s", r.URL.Path)
				}
				if r.Method != http.MethodGet {
					t.Fatalf("Expected method 'GET', got: %s", r.Method)
				}
				username, password, ok := r.BasicAuth()
				if !ok || username != "testuser" || password != "testpass" {
					t.Fatalf("Expected basic auth with username 'testuser' and password 'testpass'")
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.response))
			}))
			defer mockServer.Close()

			client := &GammuClient{
				HTTPClient: http.DefaultClient,
				Endpoint:   mockServer.URL,
				Username:   "testuser",
				Password:   "testpass",
			}

			sms, err := client.ReadSMS()
			if err != nil {
				if err != tt.wantErr {
					t.Fatalf("Expected no error, got: %v", err)
				}
			}

			// gammu-sms-gateway returns a single SMS object
			s := sms[0]
			want := tt.want

			if s.ID != want.ID {
				t.Errorf("Expected ID '%s', got: '%s'", want.ID, s.ID)
			}
			if s.Text != want.Text {
				t.Errorf("Expected Text '%s', got: '%s'", want.Text, s.Text)
			}
			if s.Number != want.Number {
				t.Errorf("Expected Number '%s', got: '%s'", want.Number, s.Number)
			}
			if s.State != want.State {
				t.Errorf("Expected State '%s', got: '%s'", want.State, s.State)
			}
			if s.Date != want.Date {
				t.Errorf("Expected Date '%s', got: '%s'", want.Date, s.Date)
			}

			validateSms := func(s Sms, expectedErr error) {
				err := s.Validate()
				if err != expectedErr {
					t.Errorf("Expected error '%v', got: '%v'", expectedErr, err)
				}
			}
			validateSms(s, tt.wantErr)
		})

	}
}
