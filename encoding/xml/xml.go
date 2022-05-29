package xml

import (
	"encoding/xml"
	"net/http"
)

type XMLEncoding struct{}

func NewXMLEncoding() *XMLEncoding {
	return &XMLEncoding{}
}

func (j *XMLEncoding) Decode(req *http.Request, v any) error {
	return xml.NewDecoder(req.Body).Decode(v)
}

func (j *XMLEncoding) Encode(rw http.ResponseWriter, v any) error {
	return xml.NewEncoder(rw).Encode(v)
}

func (j *XMLEncoding) Mimetype() string {
	return "application/xml"
}
