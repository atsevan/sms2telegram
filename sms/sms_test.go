package sms

import (
	"testing"
)

func TestSms_String(t *testing.T) {
	s := Sms{
		Date:   "2023-10-01",
		Number: "+1234567890",
		State:  "received",
		Text:   "Hello, World!",
		ID:     "1",
	}
	expected := "+1234567890 sent on 2023-10-01 (received)\nHello, World!"
	if s.String() != expected {
		t.Errorf("expected %s, got %s", expected, s.String())
	}
}

func TestSms_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sms     Sms
		wantErr error
	}{
		{
			name:    "Valid SMS",
			sms:     Sms{Date: "2023-10-01", Number: "+1234567890", State: "received", Text: "Hello, World!"},
			wantErr: nil,
		},
		{
			name:    "No Date",
			sms:     Sms{Number: "+1234567890", State: "received", Text: "Hello, World!"},
			wantErr: ErrNoDate,
		},
		{
			name:    "No Number",
			sms:     Sms{Date: "2023-10-01", State: "received", Text: "Hello, World!"},
			wantErr: ErrNoNumber,
		},
		{
			name:    "No State",
			sms:     Sms{Date: "2023-10-01", Number: "+1234567890", Text: "Hello, World!"},
			wantErr: ErrNoState,
		},
		{
			name:    "No Text",
			sms:     Sms{Date: "2023-10-01", Number: "+1234567890", State: "received"},
			wantErr: ErrNoText,
		},
		{
			name:    "No New Messages",
			sms:     Sms{},
			wantErr: ErrNoNewMessages,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.sms.Validate(); err != tt.wantErr {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
