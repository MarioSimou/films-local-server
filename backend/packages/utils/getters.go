package utils

import "fmt"

func GetSubtitleBucketKey(guid string) string {
	return fmt.Sprintf("%s/subtitle.srt", guid)
}

func GetFilmBucketKey(guid string) string {
	return fmt.Sprintf("%s/film.mp4", guid)
}

func GetImageBucketKey(guid string, ext string) string {
	return fmt.Sprintf("%s/image.%s", guid, ext)
}
