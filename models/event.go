package models

type Event struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

// TODO: after friday.. implement this