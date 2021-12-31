package models

import "time"

var (
	FilmsTableName = "films"
)

type Film struct {
	GUID        string    `json:"guid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Location    string    `json:"location"`
	Subtitle    string    `json:"subtitle"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
