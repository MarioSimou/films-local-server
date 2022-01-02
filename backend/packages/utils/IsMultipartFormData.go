package utils

import (
	"regexp"

	repoTypes "github.com/MarioSimou/songs-local-server/backend/packages/types"
)

func IsMultipartFormData(headerValue string) error {
	var regex = regexp.MustCompile("multipart\\/form-data")

	if !regex.MatchString(headerValue) {
		return repoTypes.ErrInvalidContentType
	}
	return nil
}
