package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/alert"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Timestamp string        `json:"timestamp"`
	Alerts    []alert.Alert `json:"alerts"`
	Count     int           `json:"count"`
}

// WebhookHandler sends alert notifications to an HTTP endpoint.
type WebhookHandler struct {
	URL     string
	Timeout time.Duration
	client  *http.Client
}

// NewWebhookHandler creates a WebhookHandler for the given URL.
func NewWebhookHandler(url string, timeout time.Duration) *WebhookHandler {
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &WebhookHandler{
		URL:     url,
		Timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Send posts the alerts as JSON to the configured webhook URL.
func (w *WebhookHandler) Send(alerts []alert.Alert) error {
	if len(alerts) == 0 {
		return nil
	}

	payload := WebhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Alerts:    alerts,
		Count:     len(alerts),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.URL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}

	return nil
}
