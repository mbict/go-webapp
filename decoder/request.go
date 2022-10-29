package decoder

import (
	"net/http"
)

func NewRequestDecoder(v any, tag string) (Decode, error) {
	dec, err := NewCachedDecoder(v, tag)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, v any) error {
		return dec.Decode(&RequestGetter{Request: req}, v)
	}, nil
}

type RequestGetter struct {
	*http.Request
}

func (c *RequestGetter) Get(key string) string {
	switch key {
	case `remote-addr`:
		return (*c).RemoteAddr
	case `host`:
		return (*c).Host
	case `method`:
		return (*c).Method
	case `url`:
		return (*c).URL.String()
	case `url:host`:
		return (*c).URL.Host
	case `url:query`:
		return (*c).URL.RawQuery
	case `url:path`:
		return (*c).URL.Path
	case `url:scheme`:
		return (*c).URL.Scheme
	}

	return ""
}

func (c *RequestGetter) Values(key string) []string {
	if v := c.Get(key); v != "" {
		return []string{v}
	}
	return []string{}
}
