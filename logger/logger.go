package logger

import (
	"encoding/json"
	"fmt"
	s "taskTracker/store"
	"time"
)

func ShowEvents(filename string) error {
	events, err := s.LoadItems[s.Event](filename)
	if err != nil {
		return fmt.Errorf("parsing error in showEvents: %w", err)
	}
	jsonData, err := json.MarshalIndent(events, "", "\t")
	if err != nil {
		return fmt.Errorf("converse error in showEvents")
	}
	fmt.Println(string(jsonData))
	return nil
}

func Log(fileName string, text string, err error) error {
	var errText string
	if err != nil {
		errText = err.Error()
	} else {
		errText = ""
	}
	event := s.Event{
		CreatedAT: time.Now().Format("2006-01-02 15:04"),
		UserInput: text,
		ErrorText: errText,
	}
	events, err := s.LoadItems[s.Event](fileName)
	if err != nil {
		return fmt.Errorf("parsing error in eventLogger: %w", err)
	}
	events = append(events, event)
	if err = s.SaveItems(fileName, events); err != nil {
		return fmt.Errorf("saving error in saveEvents: %w", err)
	}
	return nil
}
