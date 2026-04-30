package utils

import (
	"fmt"
	"os"
	"strings"
)

// ValidateName checks that a collection/config name is acceptable.
func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf("name contains invalid characters")
	}
	if len(name) > 255 {
		return fmt.Errorf("name is too long (max 255 characters)")
	}
	return nil
}

// ValidatePath checks that a filesystem path exists.
func ValidatePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("path cannot be empty")
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("path %q does not exist: %w", path, err)
	}
	return nil
}
