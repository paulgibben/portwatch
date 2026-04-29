package scanner

import "sync"

// PortSeverityStore evaluates and caches severity levels for ports.
type PortSeverityStore struct {
	mu     sync.RWMutex
	config PortSeverityConfig
	cache  map[string]SeverityLevel
}

// NewPortSeverityStore creates a store from the given config.
func NewPortSeverityStore(cfg PortSeverityConfig) *PortSeverityStore {
	return &PortSeverityStore{
		config: cfg,
		cache:  make(map[string]SeverityLevel),
	}
}

// Evaluate returns the severity level for a given port and protocol.
// Results are cached after first lookup.
func (s *PortSeverityStore) Evaluate(port int, protocol string) SeverityLevel {
	key := portKey(port, protocol)

	s.mu.RLock()
	if level, ok := s.cache[key]; ok {
		s.mu.RUnlock()
		return level
	}
	s.mu.RUnlock()

	level := s.resolve(port, protocol)

	s.mu.Lock()
	s.cache[key] = level
	s.mu.Unlock()

	return level
}

func (s *PortSeverityStore) resolve(port int, protocol string) SeverityLevel {
	for _, rule := range s.config.Rules {
		if rule.Port != port {
			continue
		}
		if rule.Protocol != "" && rule.Protocol != protocol {
			continue
		}
		return rule.Severity
	}
	return s.config.DefaultSeverity
}

// UpdateConfig replaces the current config and clears the cache.
func (s *PortSeverityStore) UpdateConfig(cfg PortSeverityConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
	s.cache = make(map[string]SeverityLevel)
}

func portKey(port int, protocol string) string {
	if protocol == "" {
		protocol = "tcp"
	}
	return protocol + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
