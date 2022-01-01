package utils

import "errors"

var (
	ErrImageNotFound      = errors.New("error: image not found")
	ErrS3ObjectNotFound   = errors.New("error: s3 object not found")
	ErrInvalidContentType = errors.New("error: invalid content type")
	ErrSongNotFound       = errors.New("error: song not found")
)
