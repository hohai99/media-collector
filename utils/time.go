package utils

import "fmt"

// FormatDuration converts seconds to a human-readable string (e.g. "3m 15s").
func FormatDuration(seconds int) string {
	if seconds < 0 {
		seconds = 0
	}
	h := seconds / 3600
	m := (seconds % 3600) / 60
	s := seconds % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// ParseDurationSeconds parses a simple "NmNs" or plain seconds string.
// For simplicity this just accepts an integer number of seconds.
func ParseDurationSeconds(input int) int {
	if input < 0 {
		return 0
	}
	return input
}
