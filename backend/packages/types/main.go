package types

import (
	"errors"
	"time"
)

type FileType = string
type EnvironmentType = string
type Field string

var (
	ImageType             FileType        = "image"
	SongType              FileType        = "song"
	ImageField            Field           = "image"
	LocationField         Field           = "location"
	Production            EnvironmentType = "prod"
	Dev                   EnvironmentType = "dev"
	SongsTableName                        = "songs"
	ErrImageNotFound                      = errors.New("error: image not found")
	ErrS3ObjectNotFound                   = errors.New("error: s3 object not found")
	ErrInvalidContentType                 = errors.New("error: invalid content type")
	ErrSongNotFound                       = errors.New("error: song not found")
	Pending                               = "pending"
	Done                                  = "done"
)

type Song struct {
	GUID           string    `json:"guid"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
	ImageStatus    string    `json:"imageStatus,omitempty"`
	Location       string    `json:"location"`
	LocationStatus string    `json:"locationStatus,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
