package notifier

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Type represents the kind of notifier to use.
type Type string

const (
	TypeStdout  Type = "stdout"
	TypeWebhook Type = "webhook"
	TypeCommand Type = "command"
)

// Config holds notifier configuration.
type Config struct {
	Type    Type   `json:"type"`
	Target  string `json:"target"`  // URL for webhook, command string for command
	Format  string `json:"format"`  // "text" or "json"
}

// Notifier dispatches alerts via a configured backend.
type Notifier struct {
	cfg Config
}

// New creates a new Notifier from the given config.
func New(cfg Config) *Notifier {
	if cfg.Type == "" {
		cfg.Type = TypeStdout
	}
	if cfg.Format == "" {
		cfg.Format = "text"
	}
	return &Notifier{cfg: cfg}
}

// Send dispatches the alert using the configured backend.
func (n *Notifier) Send(a *alert.Alert) error {
	switch n.cfg.Type {
	case TypeStdout:
		return n.sendStdout(a)
	case TypeCommand:
		return n.sendCommand(a)
	default:
		return fmt.Errorf("unsupported notifier type: %s", n.cfg.Type)
	}
}

func (n *Notifier) sendStdout(a *alert.Alert) error {
	for _, e := range a.Events {
		fmt.Fprintf(os.Stdout, "[portwatch] %s port %d/%s\n", e.Kind, e.Port, e.Proto)
	}
	return nil
}

func (n *Notifier) sendCommand(a *alert.Alert) error {
	if n.cfg.Target == "" {
		return fmt.Errorf("command notifier requires a target command")
	}
	parts := strings.Fields(n.cfg.Target)
	var lines []string
	for _, e := range a.Events {
		lines = append(lines, fmt.Sprintf("%s %d/%s", e.Kind, e.Port, e.Proto))
	}
	payload := strings.Join(lines, "\n")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = strings.NewReader(payload)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
