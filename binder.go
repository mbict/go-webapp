package webapp

import "github.com/mbict/webapp/decoder"

const defaultTag = "default"
const pathTag = "path"
const queryTag = "query"
const headerTag = "header"

type decoderFactory func(v any, tag string) (decoder.Decode, error)

var defaultDecodeBinders = []struct {
	tag     string
	factory decoderFactory
}{
	{
		tag:     defaultTag,
		factory: decoder.NewDefaultDecoder,
	},
	{
		tag:     headerTag,
		factory: decoder.NewHeaderDecoder,
	},
	{
		tag:     queryTag,
		factory: decoder.NewQueryDecoder,
	},
	{
		tag:     pathTag,
		factory: decoder.NewPathDecoder,
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
