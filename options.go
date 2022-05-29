package webapp

import (
	"github.com/mbict/go-webapp/container"
	"github.com/mbict/go-webapp/encoding/json"
)

type Option func(ctx *HandlerContext)

func AcceptsJson(mediatype ...string) Option {
	return func(ctx *HandlerContext) {
		ctx.decoderNegotiator.Register("application/json", json.NewJsonEncoding(), mediatype...)
	}
}

func OutputsJson() Option {
	return func(ctx *HandlerContext) {
		ctx.encoderNegotiator.Register("application/json", json.NewJsonEncoding())
	}
}

func WithContainer(container container.Container) Option {
	return func(ctx *HandlerContext) {
		//container: container,
	}
}
