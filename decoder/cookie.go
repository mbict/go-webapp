package decoder

import (
	"net/http"
)

func NewCookieDecoder(v any, tag string) (Decode, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, v any) error {
		return dec.Decode(CookieGetter(req.Cookies()), v)
	}, nil
}

type CookieGetter []*http.Cookie

func (c CookieGetter) Get(key string) string {
	for i := range c {
		if c[i].Name == key {
			return c[i].Value
		}
	}
	return ""
}

func (c CookieGetter) Values(key string) []string {
	var res []string
	for i := range c {
		if c[i].Name == key {
			res = append(res, c[i].Value)
		}
	}
	return res
}
