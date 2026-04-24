package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/alert"
)

func TestWebhookHandler_Send_Success(t *testing.T) {
	var received WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	alerts := []alert.Alert{
		{Port: 8080, Proto: "tcp", Change: "added"},
	}

	if err := h.Send(alerts); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if received.Count != 1 {
		t.Errorf("expected count 1, got %d", received.Count)
	}
	if len(received.Alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(received.Alerts))
	}
	if received.Alerts[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Alerts[0].Port)
	}
}

func TestWebhookHandler_Send_Empty(t *testing.T) {
	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	if err := h.Send(nil); err != nil {
		t.Fatalf("expected no error on empty alerts, got %v", err)
	}
	if called {
		t.Error("expected server not to be called for empty alerts")
	}
}

func TestWebhookHandler_Send_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, 5*time.Second)
	alerts := []alert.Alert{{Port: 443, Proto: "tcp", Change: "removed"}}

	if err := h.Send(alerts); err == nil {
		t.Error("expected error for non-2xx status, got nil")
	}
}

func TestNewWebhookHandler_DefaultTimeout(t *testing.T) {
	h := NewWebhookHandler("http://example.com", 0)
	if h.Timeout != 10*time.Second {
		t.Errorf("expected default timeout 10s, got %v", h.Timeout)
	}
}
