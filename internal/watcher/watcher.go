package watcher

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/scanner"
)

// Watcher periodically scans open ports and alerts on unexpected changes.
type Watcher struct {
	scanner  *scanner.Scanner
	notifier *alert.Notifier
	config   *rules.Config
	interval time.Duration
	snap     *scanner.Snapshot
	stopCh   chan struct{}
}

// New creates a new Watcher with the given scanner, notifier, config, and poll interval.
func New(s *scanner.Scanner, n *alert.Notifier, cfg *rules.Config, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  s,
		notifier: n,
		config:   cfg,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the polling loop in a background goroutine.
func (w *Watcher) Start() {
	go w.run()
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) run() {
	if err := w.tick(); err != nil {
		log.Printf("portwatch: initial scan error: %v", err)
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.tick(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		case <-w.stopCh:
			return
		}
	}
}

func (w *Watcher) tick() error {
	ports, err := w.scanner.OpenPorts()
	if err != nil {
		return err
	}

	next := scanner.NewSnapshot(ports)

	if w.snap != nil {
		diff := w.snap.Diff(next)
		if len(diff.Added) > 0 || len(diff.Removed) > 0 {
			alerts := alert.FromDiff(diff, w.config)
			for _, a := range alerts {
				w.notifier.Notify(a)
			}
		}
	}

	w.snap = next
	return nil
}
