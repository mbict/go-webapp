package decoder

import (
	"net/http"
)

// NewDefaultDecoder will set the default value of a struct value based on the tag value
// int Property `default:"123"` will set the value 123
func NewDefaultDecoder(v any, tag string) (Decode, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, v any) error {
		return dec.Decode(defaultsGetter{}, v)
	}, nil
}

type defaultsGetter struct{}

func (_ defaultsGetter) Get(key string) string {
	return key
}

func (_ defaultsGetter) Values(key string) []string {
	return []string{key}
}
