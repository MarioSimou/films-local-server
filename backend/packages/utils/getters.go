package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type FileType = string

var (
	ImageType FileType = "image"
	SongType  FileType = "song"
)

func GetBucketKey(ft FileType, songGUID, ext string) string {
	return fmt.Sprintf("%s/%s%s", songGUID, ft, ext)
}

func GetBucketKeyFromURL(href string) string {
	var u, _ = url.Parse(href)
	var oldPart = fmt.Sprintf("/%s/", S3_BUCKET_NAME)
	return strings.Replace(u.Path, oldPart, "", -1)
}
