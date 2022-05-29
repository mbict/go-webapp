package decoder

import (
	"github.com/mbict/httprouter"
	"net/http"
)

func NewPathDecoder(v any, tag string) (Decode, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, v any) error {
		params := httprouter.ParamsFromContext(req.Context())
		return dec.Decode(ParamsGetter(params), v)
	}, nil
}

type ParamsGetter []httprouter.Param

func (ps ParamsGetter) Get(key string) string {
	for i := range ps {
		if ps[i].Key == key {
			return ps[i].Value
		}
	}
	return ""
}

func (ps ParamsGetter) Values(key string) []string {
	for i := range ps {
		if ps[i].Key == key {
			return []string{ps[i].Value}
		}
	}
	return nil
}
