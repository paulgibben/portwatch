package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"portwatch/internal/rules"
	"portwatch/internal/watcher"
)

// Daemon holds the top-level runtime state for portwatch.
type Daemon struct {
	cfg     *rules.Config
	watcher *watcher.Watcher
	logger  *log.Logger
}

// New creates a new Daemon from the given config path.
func New(configPath string) (*Daemon, error) {
	cfg, err := rules.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	w := watcher.New(cfg)

	return &Daemon{
		cfg:     cfg,
		watcher: w,
		logger:  log.New(os.Stdout, "[portwatch] ", log.LstdFlags),
	}, nil
}

// Run starts the daemon and blocks until a termination signal is received.
func (d *Daemon) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	d.logger.Println("starting portwatch daemon")

	if err := d.watcher.Start(ctx); err != nil {
		return err
	}

	select {
	case sig := <-sigCh:
		d.logger.Printf("received signal %s, shutting down", sig)
		cancel()
	case <-ctx.Done():
	}

	d.watcher.Stop()
	d.logger.Println("portwatch daemon stopped")
	return nil
}
