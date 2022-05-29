package json

import (
	"github.com/goccy/go-json"
	"net/http"
)

type JsonEncoding struct{}

func NewJsonEncoding() *JsonEncoding {
	return &JsonEncoding{}
}

func (j *JsonEncoding) Decode(req *http.Request, v any) error {
	return json.NewDecoder(req.Body).Decode(v)
}

func (j *JsonEncoding) Encode(rw http.ResponseWriter, v any) error {
	return json.NewEncoder(rw).Encode(v)
}

func (j *JsonEncoding) Mimetype() string {
	return "application/json"
}
