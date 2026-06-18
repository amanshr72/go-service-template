// client.go — port + real adapter for an external notification vendor.
// The vendor sends emails; we call POST {BaseURL}/v1/send with a JSON payload.
package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// SendEmailRequest — contract: what the vendor API expects
type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SendEmailResponse — contract: what the vendor API returns
type SendEmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// Client is the port — business code depends on this, never on http.Client directly.
// This is what makes swapping real vendor URL <-> mock server URL trivial.
type Client interface {
	SendEmail(req SendEmailRequest) (*SendEmailResponse, error)
}

// httpClient is the real adapter — talks to whatever BaseURL points to.
// Same struct hits the real vendor in prod and our mock server in tests/dev.
type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient — baseURL is injected, not hardcoded.
// Prod: https://api.realvendor.com
// Dev/test: http://localhost:8089 (our mock server)
func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // never let an external call hang forever
		},
	}
}

func (c *httpClient) SendEmail(req SendEmailRequest) (*SendEmailResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/send", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call vendor: %w", err) // network error, timeout etc.
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("vendor returned non-200 status")
	}

	var out SendEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &out, nil
}
