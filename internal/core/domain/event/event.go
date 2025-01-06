package event

import "fmt"

type Event interface {
	AsMessage() string
}

type FailedUpdateEvent struct {
	Message string
	Error   error
}

func NewFailedUpdateEvent(message string, err error) *FailedUpdateEvent {
	return &FailedUpdateEvent{
		Message: message,
		Error:   err,
	}
}

func (e FailedUpdateEvent) AsMessage() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Error)
}

type ChangeEvent struct {
	Message string
}

func NewChangeEvent(message string) *ChangeEvent {
	return &ChangeEvent{
		Message: message,
	}
}

func (e ChangeEvent) AsMessage() string {
	return e.Message
}
