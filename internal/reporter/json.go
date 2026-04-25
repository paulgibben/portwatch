package reporter

import (
	"encoding/json"
	"fmt"
	"io"
)

// writeJSON serialises rep as indented JSON to w.
func writeJSON(w io.Writer, rep Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rep); err != nil {
		return fmt.Errorf("reporter: json encode: %w", err)
	}
	return nil
}

// ParseReport deserialises a JSON-encoded Report from r.
func ParseReport(r io.Reader) (Report, error) {
	var rep Report
	if err := json.NewDecoder(r).Decode(&rep); err != nil {
		return Report{}, fmt.Errorf("reporter: json decode: %w", err)
	}
	return rep, nil
}
