package models

import (
	"os"
	"time"
)

type TableName = string
type BucketName = string

var (
	SongsTableName    TableName  = os.Getenv("AWS_DYNAMODB_TABLE_NAME")
	ServiceBucketName BucketName = os.Getenv("AWS_S3_BUCKET_NAME")
)

type Song struct {
	GUID        string    `json:"guid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Href        string    `json:"href"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PostSong struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"max=1000"`
	GUID        string    `json:"-"`
	Image       string    `json:"-"`
	Href        string    `json:"-"`
	CreatedAt   time.Time `json:"-"`
}

type PutSong struct {
	Name        string `json:"name" validate:"required_without_all=Description"`
	Description string `json:"description" validate:"max=1000"`
}
