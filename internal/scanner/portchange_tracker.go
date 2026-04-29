package scanner

// PortChangeTracker builds a PortChangeLog from a Diff.
type PortChangeTracker struct {
	log *PortChangeLog
}

// NewPortChangeTracker returns a tracker backed by the given log.
func NewPortChangeTracker(log *PortChangeLog) *PortChangeTracker {
	if log == nil {
		log = NewPortChangeLog()
	}
	return &PortChangeTracker{log: log}
}

// Log returns the underlying PortChangeLog.
func (t *PortChangeTracker) Log() *PortChangeLog {
	return t.log
}

// Apply records all changes from d into the log.
func (t *PortChangeTracker) Apply(d Diff) {
	for _, p := range d.Added {
		t.log.Record(p.Port, p.Proto, ChangeAdded)
	}
	for _, p := range d.Removed {
		t.log.Record(p.Port, p.Proto, ChangeRemoved)
	}
}

// Since returns all entries detected at or after the given Unix timestamp (seconds).
func (t *PortChangeTracker) Since(unixSec int64) []PortChange {
	var out []PortChange
	for _, e := range t.log.Entries {
		if e.DetectedAt.Unix() >= unixSec {
			out = append(out, e)
		}
	}
	return out
}

// CountByType returns how many entries match the given ChangeType.
func (t *PortChangeTracker) CountByType(ct ChangeType) int {
	n := 0
	for _, e := range t.log.Entries {
		if e.Change == ct {
			n++
		}
	}
	return n
}
