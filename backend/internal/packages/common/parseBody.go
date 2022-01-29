package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
)

func ParseBody(body string, model interface{}) error {
	var validate = validator.New()

	if e := json.NewDecoder(strings.NewReader(body)).Decode(model); e != nil {
		return e
	}
	if e := validate.Struct(model); e != nil {
		return e
	}

	return nil
}

func GetGUIDFromParameters(parameters map[string]string) (*string, error) {
	var validate = validator.New()
	if guid, ok := parameters["guid"]; ok {
		if e := validate.Var(guid, "required,uuid4"); e != nil {
			return nil, e
		}
		return &guid, nil
	}
	return nil, fmt.Errorf("error: guid not found")
}

func GetLimitFromQueryParameters(parameters map[string]string) (*int32, error) {
	var limitStr string
	var limit int64
	var defaultLimit int32 = 20
	var ok bool
	var e error

	if limitStr, ok = parameters["limit"]; ok {
		if limit, e = strconv.ParseInt(limitStr, 10, 32); e != nil {
			return nil, e
		}
		var limitCasted = int32(limit)
		return &limitCasted, nil
	}
	return &defaultLimit, nil
}

type MultipartResponse struct {
	Filename  string
	Body      []byte
	MediaType string
	Ext       string
}

func ParseMultipartContentType(headers map[string]string, body string) (*MultipartResponse, error) {
	var contentType string
	var ok bool
	var mediaType string
	var mediaParams map[string]string
	var e error

	if contentType, ok = headers["Content-Type"]; !ok {
		return nil, fmt.Errorf("error: Unable to find content-type")
	}
	if ok := strings.Contains(contentType, "multipart/form-data"); !ok {
		return nil, fmt.Errorf("error: Invalid multipart content type")
	}

	if mediaType, mediaParams, e = mime.ParseMediaType(contentType); e != nil {
		return nil, e
	}

	var multipartReader = multipart.NewReader(strings.NewReader(body), mediaParams["boundary"])

	var bf = bytes.NewBuffer(nil)
	var iteration = 0
	var filename string
	for {
		var part *multipart.Part
		var partBf []byte
		var e error

		if part, e = multipartReader.NextPart(); e != nil {
			if e == io.EOF {
				break
			}
			return nil, e
		}
		defer part.Close()

		if iteration == 0 {
			filename = part.FileName()
		}

		if partBf, e = ioutil.ReadAll(part); e != nil {
			return nil, e
		}

		if _, e := bf.Write(partBf); e != nil {
			return nil, e
		}
	}

	return &MultipartResponse{
		Filename:  filename,
		Body:      bf.Bytes(),
		MediaType: mediaType,
		Ext:       filepath.Ext(filename),
	}, nil
}
