package types

import (
	"errors"
	"time"
)

type FileType = string
type EnvironmentType = string

var (
	ImageType             FileType        = "image"
	SongType              FileType        = "song"
	Production            EnvironmentType = "prod"
	Dev                   EnvironmentType = "dev"
	SongsTableName                        = "songs"
	ErrImageNotFound                      = errors.New("error: image not found")
	ErrS3ObjectNotFound                   = errors.New("error: s3 object not found")
	ErrInvalidContentType                 = errors.New("error: invalid content type")
	ErrSongNotFound                       = errors.New("error: song not found")
)

type Song struct {
	GUID        string    `json:"guid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
