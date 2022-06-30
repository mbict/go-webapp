package webapp

import (
	"github.com/mbict/go-webapp/container"
	"github.com/mbict/go-webapp/encoding/json"
	"github.com/mbict/go-webapp/encoding/xml"
)

var DefaultOptions = Options{
	WithContainer(container.Default),
	AcceptsJson("*/*"),
	OutputsJson("*/*"),
	WithDefaultJSONOutputEncoding(),
	//WithErrorHandler(func(err error) error {
	//	return err
	//}),
}

type Options []Option

func (o *Options) Add(option ...Option) Options {
	return append(*o, option...)
}

type Option func(ctx *HandlerContext)

func AcceptsJson(mediatypeAlias ...string) Option {
	return func(ctx *HandlerContext) {
		ctx.decoderNegotiator.Register("application/json", json.NewJsonEncoding(), mediatypeAlias...)
	}
}

func OutputsJson(mediatypeAlias ...string) Option {
	return func(ctx *HandlerContext) {
		ctx.encoderNegotiator.Register("application/json", json.NewJsonEncoding(), mediatypeAlias...)
	}
}

func WithDefaultOutputEncoding(mimetype string) Option {
	return func(ctx *HandlerContext) {
		ctx.defaultEncoding = mimetype
	}
}

func WithDefaultJSONOutputEncoding() Option {
	return func(ctx *HandlerContext) {
		ctx.defaultEncoding = "application/json"
	}
}

func AcceptsXML(mediatype ...string) Option {
	return func(ctx *HandlerContext) {
		ctx.decoderNegotiator.Register("application/xml", xml.NewXMLEncoding(), mediatype...)
	}
}

func OutputsXML() Option {
	return func(ctx *HandlerContext) {
		ctx.encoderNegotiator.Register("application/xml", xml.NewXMLEncoding())
	}
}

func WithContainer(container container.Container) Option {
	return func(ctx *HandlerContext) {
		ctx.container = container
	}
}

func WithErrorHandler(handler func(error) error) Option {
	return func(ctx *HandlerContext) {
		ctx.errorHandler = handler
	}
}
