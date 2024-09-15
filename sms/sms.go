package sms

import (
	"fmt"
)

var (
	ErrNoNewMessages = fmt.Errorf("no new messages")
	ErrNoDate        = fmt.Errorf("date is required")
	ErrNoNumber      = fmt.Errorf("number is required")
	ErrNoState       = fmt.Errorf("state is required")
	ErrNoText        = fmt.Errorf("text is required")
)

// Sms represents a SMS message.
type Sms struct {
	Date   string `json:"Date"`
	Number string `json:"number"`
	State  string `json:"State"`
	Text   string `json:"Text"`
	ID     string `json:"ID"`
}

func (s Sms) String() string {
	return fmt.Sprintf("%s sent on %s (%s)\n%s", s.Number, s.Date, s.State, s.Text)
}

// Validate validates the Sms struct.
func (s Sms) Validate() error {
	// if the SMS is empty assuming that there are no new messages
	if s.Date == "" && s.Number == "" && s.State == "" && s.Text == "" {
		return ErrNoNewMessages
	}

	if s.Date == "" {
		return ErrNoDate
	}
	if s.Number == "" {
		return ErrNoNumber
	}
	if s.State == "" {
		return ErrNoState
	}
	if s.Text == "" {
		return ErrNoText
	}
	return nil
}
