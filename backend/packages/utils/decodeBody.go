package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

func Base64Decode(body string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(body)
}

func DecodeEventBody(body string, item interface{}, isBase64Encoded bool) error {
	var e error
	var bf []byte

	if isBase64Encoded {
		if bf, e = Base64Decode(body); e != nil {
			return e
		}
		body = string(bf)
	}

	if e := json.NewDecoder(strings.NewReader(body)).Decode(item); e != nil {
		return e
	}
	return nil
}
