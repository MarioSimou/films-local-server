package utils

import "errors"

var (
	ErrSubtitleNotFound   = errors.New("error: subtitle not found")
	ErrS3ObjectNotFound   = errors.New("error: s3 object not found")
	ErrInvalidContentType = errors.New("error: inalid content type")
	ErrFilmNotFound       = errors.New("error: film not found")
)
