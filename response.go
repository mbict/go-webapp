package webapp

import (
	"encoding/xml"
	"net/http"
)

type CreatedResponse struct {
	empty
	url string
}

func (c *CreatedResponse) StatusCode() int {
	return http.StatusCreated
}

func (c *CreatedResponse) Header() http.Header {
	h := http.Header{}
	h.Add("Location", c.url)
	return h
}

func NewCreatedResponse(url string) *CreatedResponse {
	return &CreatedResponse{
		empty: EmptyResponse,
		url:   url,
	}
}

type Empty struct{}

var EmptyResponse = empty(false)

type empty bool

var emptyBytes = []byte(``)

func (r empty) MarshalText() ([]byte, error) {
	return emptyBytes, nil
}

func (r empty) MarshalJSON() ([]byte, error) {
	return emptyBytes, nil
}

func (r empty) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	return nil
}

func (r empty) MarshalYAML() (interface{}, error) {
	return emptyBytes, nil
}
