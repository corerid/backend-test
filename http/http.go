package http

import (
	"encoding/json"
	"io"
	"net/http"
)

//go:generate mockgen -source=http.go -destination=mocks/http.go

const (
	MethodPost = http.MethodPost
)

type Client interface {
	Send(method string, url string, request *Request) (*http.Response, error)
}

type HTTPClient interface {
	// Do function to execute
	Do(req *http.Request) (*http.Response, error)
}

type ClientWrapper struct {
	client HTTPClient
}

func New(client HTTPClient) *ClientWrapper {
	return &ClientWrapper{client: client}
}

func (c *ClientWrapper) Send(method string, url string, request *Request) (*http.Response, error) {
	req, err := http.NewRequest(method, url, request.Body)
	if err != nil {

	}
	req.Header.Set("Content-Type", request.ContentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func ParseHTTPResponseToStruct(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, target)
	if err != nil {
		return err
	}

	return nil
}
