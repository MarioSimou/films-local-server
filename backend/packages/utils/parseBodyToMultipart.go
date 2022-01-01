package utils

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
)

type Multipart struct {
	Type     string
	Params   map[string]string
	Reader   *multipart.Reader
	Filename string
	Ext      string
	Body     []byte
}

func ParseBodyToMultipart(event events.APIGatewayProxyRequest) (*Multipart, error) {
	var mediaType string
	var mediaTypeParams map[string]string
	var e error
	var contentTypeHeader = event.Headers["Content-Type"]
	var body []byte
	var part *multipart.Part

	if mediaType, mediaTypeParams, e = mime.ParseMediaType(contentTypeHeader); e != nil {
		return nil, e
	}

	if body, e = Base64Decode(event.Body); e != nil {
		return nil, e
	}

	var r = multipart.NewReader(bytes.NewReader(body), mediaTypeParams["boundary"])

	if part, e = r.NextPart(); e != nil {
		return nil, e
	}

	if body, e = ioutil.ReadAll(part); e != nil {
		return nil, e
	}

	return &Multipart{
		Type:     mediaType,
		Params:   mediaTypeParams,
		Reader:   r,
		Body:     body,
		Filename: part.FileName(),
		Ext:      filepath.Ext(part.FileName()),
	}, nil
}
