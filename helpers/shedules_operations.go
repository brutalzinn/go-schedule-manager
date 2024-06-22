package helpers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/brutalzinn/go-schedule-manager/models"
	"github.com/google/uuid"
)

func ReadSchedules(day string) ([]models.Schedule, error) {
	jsonFile, err := os.ReadFile(fmt.Sprintf("static/routines/%s.json", day))
	if err != nil {
		return nil, fmt.Errorf("could not read JSON file for day %s: %v", day, err)
	}
	var schedules []models.Schedule
	err = json.Unmarshal(jsonFile, &schedules)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON for day %s: %v", day, err)
	}

	return schedules, nil
}

func ProcessCSV(file io.Reader) (map[string][]models.Schedule, error) {
	reader := csv.NewReader(file)
	reader.Comma = ','

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("could not read the header row: %v", err)
	}

	schedulesByDay := make(map[string][]models.Schedule)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("could not read the row: %v", err)
		}

		timeRange := record[0]
		for i, day := range header[1:] {
			message := record[i+1]
			if message != "" {
				timeParts := strings.Split(timeRange, " - ")
				schedule := models.Schedule{
					ID:      uuid.New().String(),
					Time:    fmt.Sprintf("%s:00", timeParts[0]),
					Content: message,
					UseTTS:  true,
				}
				schedulesByDay[day] = append(schedulesByDay[day], schedule)
			}
		}
	}
	return schedulesByDay, nil
}
