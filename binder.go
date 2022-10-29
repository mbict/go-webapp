package webapp

import "github.com/mbict/go-webapp/decoder"

const defaultTag = "default"
const pathTag = "path"
const queryTag = "query"
const headerTag = "header"
const cookieTag = "cookie"
const requestTag = "request"

type decoderFactory func(v any, tag string) (decoder.Decode, error)

var defaultDecodeBinders = []struct {
	tag     string
	factory decoderFactory
}{
	{
		tag:     headerTag,
		factory: decoder.NewHeaderDecoder,
	},
	{
		tag:     queryTag,
		factory: decoder.NewQueryDecoder,
	},
	{
		tag:     cookieTag,
		factory: decoder.NewCookieDecoder,
	},
	{
		tag:     pathTag,
		factory: decoder.NewPathDecoder,
	},
	{
		tag:     requestTag,
		factory: decoder.NewRequestDecoder,
	},
}

func BuildArgumentsBinder(v any) (Decoder, error) {
	res := decoders{}
	for _, c := range defaultDecodeBinders {
		if hasTag(v, c.tag) {
			dec, err := c.factory(v, c.tag)
			if err != nil {
				return nil, err
			}
			res = append(res, dec)
		}
	}
	return &res, nil
}
