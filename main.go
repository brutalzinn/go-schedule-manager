package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/brutalzinn/go-schedule-manager/helpers"
	"github.com/brutalzinn/go-schedule-manager/models"
	"github.com/brutalzinn/go-schedule-manager/translator"
	"github.com/brutalzinn/go-schedule-manager/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var (
	clients         = make(map[*websocket.Conn]bool)
	schedulesOfDay  []models.Schedule
	currentDay      = utils.GetCurrentDay()
	translateBundle translator.TranslatorBundle
)

func main() {
	translateBundle = translator.New(language.Portuguese)
	loadSchedulesOfDay(currentDay)
	app := fiber.New()
	app.Static("/", "./static")
	app.Post("/upload", uploadHandler)
	app.Get("/json/:day", jsonHandler)
	app.Post("/save/:day", saveHandler)
	app.Post("/create", createScheduleHandler)
	app.Get("/refresh/:day", func(c *fiber.Ctx) error {
		day := translateBundle.T(c.Params("day"))
		loadSchedulesOfDay(day)
		initialEvents := GetInitialEvents(day)
		notifyClients(initialEvents)
		return c.SendString("Schedules updated successfully")
	})

	app.Get("/ws/:day", websocket.New(wsHandler))
	_, err := helpers.CreateAudioTTS("Testando idioma português")
	if err != nil {
		logrus.Error(err)
	}
	go scheduleChecker()
	logrus.Println("Server started at :8000")
	logrus.Info(app.Listen(":8000"))
}

func scheduleChecker() {
	for {
		now := time.Now().Format("15:04:05")
		for i, schedule := range schedulesOfDay {
			_, err := time.Parse("15:04:05", schedule.Time)
			if err != nil {
				logrus.Printf("Failed to parse time: %v", err)
				continue
			}
			if now == schedule.Time {
				logrus.Info("Executing schedule")
				triggerSchedule(schedule)
				logrus.Info("Removing from schedule")
				schedulesOfDay = append(schedulesOfDay[:i], schedulesOfDay[i+1:]...)
				notifyClients(GetInitialEvents(currentDay))
			}
		}
		logrus.Info("current time:", now)
		time.Sleep(1 * time.Second)
	}
}

func triggerSchedule(schedule models.Schedule) {
	log.Printf("Executing schedule: %s", schedule.Content)
	if schedule.UseTTS {
		audioFileName, err := helpers.CreateAudioTTS(schedule.Content)
		if err != nil {
			logrus.Error("Audio file created: ", audioFileName)
		}
	}
}

func createScheduleHandler(c *fiber.Ctx) error {
	content := c.FormValue("content")
	time := c.FormValue("time")
	useTTS := c.FormValue("useTTS") == "on" // Checkbox value check
	if !isValidTimeFormat(time) {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid time format. Please use HH:MM.")
	}
	///TODO: implement http request when create new shedule. HO can inform to alexa device that this is triggered too
	day := translateBundle.T(utils.GetCurrentDay())
	schedule := models.Schedule{
		ID:      "",
		Time:    time,
		Content: content,
		UseTTS:  useTTS,
	}
	jsonData, err := json.MarshalIndent(schedule, "", "    ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error converting to JSON")
	}
	err = os.WriteFile(fmt.Sprintf("static/routines/%s.json", day), jsonData, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving JSON file")
	}
	return c.SendString("New schedule created successfully")
}

func loadSchedulesOfDay(day string) error {
	schedules, err := helpers.ReadSchedules(translateBundle.T(day))
	if err != nil {
		return err
	}
	schedulesOfDay = schedules
	return nil
}

func uploadHandler(c *fiber.Ctx) error {
	file, err := c.FormFile("csvfile")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error reading file")
	}
	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening file")
	}
	defer fileHeader.Close()
	schedulesByDay, err := helpers.ProcessCSV(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error processing CSV file")
	}
	for day, schedules := range schedulesByDay {
		jsonData, err := json.MarshalIndent(schedules, "", "    ")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error converting to JSON")
		}
		err = os.WriteFile(fmt.Sprintf("static/routines/%s.json", day), jsonData, 0644)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error saving JSON file")
		}
	}
	return c.SendString("CSV file processed and JSON files created with success")
}

func jsonHandler(c *fiber.Ctx) error {
	day := translateBundle.T(c.Params("day"))
	jsonFile := fmt.Sprintf("static/routines/%s.json", day)
	return c.SendFile(jsonFile)
}

func saveHandler(c *fiber.Ctx) error {
	day := translateBundle.T(c.Params("day"))
	var schedules []models.Schedule
	if err := c.BodyParser(&schedules); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error decoding JSON")
	}
	jsonData, err := json.MarshalIndent(schedules, "", "    ")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error converting to JSON")
	}
	err = os.WriteFile(fmt.Sprintf("static/routines/%s.json", day), jsonData, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving JSON file")
	}
	err = notifyClients(schedules)
	if err != nil {
		logrus.Printf("Error notifying clients for day %s: %v", day, err)
	}
	return c.SendString("JSON data has been saved")
}

func wsHandler(c *websocket.Conn) {
	day := translateBundle.T(c.Params("day"))
	logrus.Printf("WebSocket connection established for day %s", day)
	clients[c] = true
	defer func() {
		delete(clients, c)
		c.Close()
		logrus.Printf("WebSocket connection closed for day %s", day)
	}()
	initialEvent := GetInitialEvents(day)
	if err := c.WriteJSON(initialEvent); err != nil {
		logrus.Printf("Error sending initial events to client: %v", err)
	}
	for {
		var msg models.Schedule
		err := c.ReadJSON(&msg)
		if err != nil {
			logrus.Println("WebSocket connection closed by client")
			break
		}
		logrus.Printf("Received message from client: %+v\n", msg)
	}
}

func GetInitialEvents(day string) models.Event {
	nextSchedules := make([]any, 0)
	historySchedules := make([]any, 0)
	schedules, err := helpers.ReadSchedules(day)
	if err != nil {
		logrus.Printf("Error reading schedules for day %s: %v", day, err)
	}
	for _, schedule := range schedules {
		if schedule.Time > time.Now().Format("15:04:05") {
			nextSchedules = append(nextSchedules, schedule.ToEvent())
		}
		if schedule.Time < time.Now().Format("15:04:05") {
			historySchedules = append(historySchedules, schedule.ToEvent())
		}
	}
	event := models.Event{
		Action: "initialEvents",
		Data: map[string]any{
			"next":    nextSchedules,
			"history": historySchedules,
		},
	}
	return event
}

func notifyClients(data any) error {
	for client := range clients {
		// if client.Params("day") == day {
		if err := client.WriteJSON(data); err != nil {
			logrus.Printf("Error sending WebSocket message to client: %v", err)
		}
		// }
	}
	return nil
}

func isValidTimeFormat(scheduleTime string) bool {
	_, err := time.Parse("15:04:05", scheduleTime)
	return err == nil
}
