package http

import (
	"bytes"
	"io"
)

type Request struct {
	Body        io.Reader
	ContentType string
}

func NewRequest(body []byte, contentType string) *Request {
	return &Request{
		Body:        bytes.NewReader(body),
		ContentType: contentType,
	}
}
