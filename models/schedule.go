package models

import (
	"fmt"

	"github.com/brutalzinn/go-schedule-manager/utils"
)

type Schedule struct {
	ID      string `json:"id"`
	Time    string `json:"time"`
	Content string `json:"content"`
	UseTTS  bool   `json:"useTTS"`
}

func (schedule Schedule) ToEvent() Event {
	if schedule.UseTTS {
		audioFilename := utils.GetMD5Hash(schedule.Content)
		event := Event{
			Action: "audio",
			Data: map[string]any{
				"time":  schedule.Time,
				"audio": fmt.Sprintf("audio/%s.mp3", audioFilename),
				"html":  fmt.Sprintf("<audio controls><source src='audio/%s.mp3' type='audio/mp3'></audio> <br/> Trigget at: %s <br/> Text: %s", audioFilename, schedule.Time, schedule.Content),
			},
		}
		return event
	}
	///lazy way
	return Event{Action: "NOTHING"}
}
