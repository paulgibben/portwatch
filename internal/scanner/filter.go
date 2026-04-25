package scanner

import "github.com/user/portwatch/internal/rules"

// Filter holds the configuration for filtering scanned ports.
type Filter struct {
	config *rules.Config
}

// NewFilter creates a new Filter using the provided rules config.
func NewFilter(cfg *rules.Config) *Filter {
	return &Filter{config: cfg}
}

// Apply takes a list of open ports and returns only those that should
// be monitored, respecting the configured port ranges and ignore lists.
func (f *Filter) Apply(ports []int) []int {
	if f.config == nil {
		return ports
	}

	var result []int
	for _, p := range ports {
		if f.shouldIgnore(p) {
			continue
		}
		if f.inRange(p) {
			result = append(result, p)
		}
	}
	return result
}

// shouldIgnore returns true if the port is in the ignore list.
func (f *Filter) shouldIgnore(port int) bool {
	for _, ignored := range f.config.IgnorePorts {
		if ignored == port {
			return true
		}
	}
	return false
}

// inRange returns true if the port falls within any configured range.
// If no ranges are configured, all ports are considered in range.
func (f *Filter) inRange(port int) bool {
	if len(f.config.PortRanges) == 0 {
		return true
	}
	for _, r := range f.config.PortRanges {
		if port >= r.Start && port <= r.End {
			return true
		}
	}
	return false
}
