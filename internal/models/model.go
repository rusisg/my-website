package models

import "time"

type Note struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Paragraph   string    `json:"paragraph"`
	TimeCreated time.Time `json:"time_created"`
	Posted      bool      `json:"posted"`
}
