package webapp

import (
	"github.com/mbict/webapp/decoder"
	"net/http"
)

type Decoder interface {
	Decode(req *http.Request, v any) error
}

type decoders []decoder.Decode

func (d *decoders) Decode(req *http.Request, v any) error {
	for _, decodeFunc := range *d {
		if err := decodeFunc(req, v); err != nil {
			return err
		}
	}
	return nil
}
