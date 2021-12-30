package utils

import (
	"encoding/json"
	"strings"
)

func DecodeEventBody(body string, item interface{}) error {
	if e := json.NewDecoder(strings.NewReader(body)).Decode(item); e != nil {
		return e
	}
	return nil
}
