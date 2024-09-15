package sms

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GammuClient struct {
	HTTPClient *http.Client
	Endpoint   string
	Username   string
	Password   string
}

func (g *GammuClient) ReadSMS() ([]Sms, error) {
	req, err := http.NewRequest(http.MethodGet, g.Endpoint+"/getsms", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(g.Username, g.Password)

	resp, err := g.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sms Sms
	err = json.NewDecoder(resp.Body).Decode(&sms)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return []Sms{sms}, sms.Validate()
}
