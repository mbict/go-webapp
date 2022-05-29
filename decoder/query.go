package decoder

import (
	"net/http"
)

func NewQueryDecoder(v any, tag string) (Decode, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, v any) error {
		return dec.Decode(MapGetter(req.URL.Query()), v)
	}, nil
}
