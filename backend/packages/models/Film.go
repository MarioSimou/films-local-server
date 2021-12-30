package models

import "time"

var (
	FilmsTableName = "films"
)

type Film struct {
	GUID        string    `json:"guid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
