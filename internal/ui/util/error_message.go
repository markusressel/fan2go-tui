package util

import (
	"errors"
	"fan2go-tui/internal/client"
)

// ErrorMessage returns the text shown in the UI for an error of the API client. Go error strings start lowercase,
// so the errors of the client that users see often get a capitalized message with a hint here.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, client.ErrNotFound) {
		return "Cannot reach fan2go daemon, did you enable its API?"
	}
	return err.Error()
}
