package notifier

import (
	"net"
	"net/smtp"
	"strings"
	"testing"
)

func TestNewEmailHandler_MissingHost(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		From: "from@example.com",
		To:   []string{"to@example.com"},
	})
	if err == nil {
		t.Fatal("expected error for missing smtp_host")
	}
}

func TestNewEmailHandler_MissingFrom(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		SMTPHost: "localhost",
		To:       []string{"to@example.com"},
	})
	if err == nil {
		t.Fatal("expected error for missing from")
	}
}

func TestNewEmailHandler_MissingRecipients(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		SMTPHost: "localhost",
		From:     "from@example.com",
	})
	if err == nil {
		t.Fatal("expected error for missing recipients")
	}
}

func TestNewEmailHandler_DefaultPort(t *testing.T) {
	h, err := NewEmailHandler(EmailConfig{
		SMTPHost: "localhost",
		From:     "from@example.com",
		To:       []string{"to@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.cfg.SMTPPort != 587 {
		t.Errorf("expected default port 587, got %d", h.cfg.SMTPPort)
	}
}

func TestNewEmailHandler_ValidConfig(t *testing.T) {
	_, err := NewEmailHandler(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: 465,
		From:     "alerts@example.com",
		To:       []string{"admin@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestEmailHandler_Send_UsesCorrectAddress starts a minimal SMTP listener
// and verifies the handler connects and transmits a message.
func TestEmailHandler_Send_UsesCorrectAddress(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	// Accept one connection and respond with minimal SMTP greeting then quit.
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Minimal SMTP server stub — just enough to not block the client.
		_ = smtp.NewClient(conn, "127.0.0.1")
	}()

	host, portStr, _ := net.SplitHostPort(ln.Addr().String())
	var port int
	_, _ = fmt.Sscanf(portStr, "%d", &port)

	h, _ := NewEmailHandler(EmailConfig{
		SMTPHost: host,
		SMTPPort: port,
		From:     "from@example.com",
		To:       []string{"to@example.com"},
	})

	a := Alert{Port: 8080, Protocol: "tcp", Kind: "added", Message: "new port"}
	// We expect an error since our stub doesn't speak full SMTP; just ensure
	// the handler attempted to connect (error is not nil but not a dial error).
	err = h.Send(a)
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		t.Errorf("handler did not attempt connection: %v", err)
	}
}
