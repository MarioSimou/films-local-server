package utils

import "regexp"

func IsMultipartFormData(headerValue string) error {
	var regex = regexp.MustCompile("multipart\\/form-data")

	if !regex.MatchString(headerValue) {
		return ErrInvalidContentType
	}
	return nil
}
