package webapp

import (
	"encoding/xml"
	"net/http"
)

type CreatedResponse struct {
	Empty
	url string
}

func (c CreatedResponse) StatusCode() int {
	return http.StatusCreated
}

func (c CreatedResponse) Header() http.Header {
	h := http.Header{}
	h.Add("Location", c.url)
	return h
}

func NewCreatedResponse(url string) CreatedResponse {
	return CreatedResponse{
		url: url,
	}
}

type empty interface {
	emptyResponse() bool
}

type Empty struct{}

var EmptyResponse = Empty{}

func (r Empty) emptyResponse() bool {
	return true
}

func (_ *Empty) StatusCode() int {
	return http.StatusNoContent
}

var emptyBytes = []byte(``)

func (r Empty) MarshalText() ([]byte, error) {
	return emptyBytes, nil
}

func (r Empty) MarshalJSON() ([]byte, error) {
	return emptyBytes, nil
}

func (r Empty) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	return nil
}

func (r Empty) MarshalYAML() (interface{}, error) {
	return emptyBytes, nil
}
