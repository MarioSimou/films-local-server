package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

func Base64Decode(body string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(body)
}

func DecodeEventBody(body string, item interface{}) error {
	var bf []byte
	var e error

	if bf, e = Base64Decode(body); e != nil {
		return e
	}

	if e := json.NewDecoder(bytes.NewReader(bf)).Decode(item); e != nil {
		return e
	}
	return nil
}
