package store

import (
	"time"
)

const (
	TaskJson   = "tasks.json"
	Charset    = "abcdefghijklmnopqrstuvwxyz0123456789"
	EventsJson = "events.json"
)

type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAT   time.Time `json:"created_at"`
	CompleteAT  time.Time `json:"update_at"`
}

type Event struct {
	CreatedAT string `json:"created_at"`
	UserInput string `json:"user_input"`
	ErrorText string `json:"error_text"`
}
