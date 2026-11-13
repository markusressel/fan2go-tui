package util

import (
	"errors"
	"fan2go-tui/internal/client"
	"fmt"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"nil", nil, ""},
		{"not found", client.ErrNotFound, "Cannot reach fan2go daemon, did you enable its API?"},
		{"wrapped not found", fmt.Errorf("get fans: %w", client.ErrNotFound), "Cannot reach fan2go daemon, did you enable its API?"},
		{"other", errors.New("connection refused"), "connection refused"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorMessage(tt.err); got != tt.expected {
				t.Errorf("ErrorMessage(%v) = %q, expected %q", tt.err, got, tt.expected)
			}
		})
	}
}
