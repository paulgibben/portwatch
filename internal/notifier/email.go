package notifier

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds configuration for the email notifier.
type EmailConfig struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	To       []string `json:"to"`
}

// EmailHandler sends alert notifications via email.
type EmailHandler struct {
	cfg EmailConfig
}

// NewEmailHandler creates a new EmailHandler with the given config.
func NewEmailHandler(cfg EmailConfig) (*EmailHandler, error) {
	if cfg.SMTPHost == "" {
		return nil, fmt.Errorf("email notifier: smtp_host is required")
	}
	if len(cfg.To) == 0 {
		return nil, fmt.Errorf("email notifier: at least one recipient is required")
	}
	if cfg.From == "" {
		return nil, fmt.Errorf("email notifier: from address is required")
	}
	port := cfg.SMTPPort
	if port == 0 {
		port = 587
		cfg.SMTPPort = port
	}
	return &EmailHandler{cfg: cfg}, nil
}

// Send delivers the alert message via SMTP.
func (e *EmailHandler) Send(a Alert) error {
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)

	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.SMTPHost)
	}

	subject := fmt.Sprintf("[portwatch] %s on port %d", a.Kind, a.Port)
	body := fmt.Sprintf("Port: %d\nProtocol: %s\nEvent: %s\nMessage: %s",
		a.Port, a.Protocol, a.Kind, a.Message)

	msg := []byte("To: " + strings.Join(e.cfg.To, ",") + "\r\n" +
		"From: " + e.cfg.From + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, msg)
}
