package status_message

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	StatusMessageDurationInfinite = -1 * time.Second
)

type StatusMessage struct {
	Message  string
	Duration time.Duration
	Color    tcell.Color
}

func newStatusMessage(message string) *StatusMessage {
	return &StatusMessage{
		Message:  message,
		Duration: StatusMessageDurationInfinite,
		Color:    tcell.ColorLightGray,
	}
}

func (statusMessage *StatusMessage) SetDuration(duration time.Duration) *StatusMessage {
	statusMessage.Duration = duration
	return statusMessage
}

func (statusMessage *StatusMessage) SetColor(color tcell.Color) *StatusMessage {
	statusMessage.Color = color
	return statusMessage
}

func NewErrorStatusMessage(message string) *StatusMessage {
	return newStatusMessage(message).SetColor(tcell.ColorRed)
}

func NewWarningStatusMessage(message string) *StatusMessage {
	return newStatusMessage(message).SetColor(tcell.ColorYellow)
}

func NewInfoStatusMessage(message string) *StatusMessage {
	return newStatusMessage(message)
}
