package webapp

//https://github.com/abemedia/go-don

import (
	"context"
	"github.com/mbict/go-webapp/container"
	"github.com/mbict/go-webapp/encoding/json"
	"log"
	"net/http"
)

// StatusCoder allows you to customise the HTTP response code.
type StatusCoder interface {
	StatusCode() int
}

// Headerer allows you to customise the HTTP headers.
type Headerer interface {
	Header() http.Header
}

// Handle is the type for the http headers
type Handle[T any, O any] func(ctx context.Context, request T) (O, error)

//type HttpHandle func(*fasthttp.RequestCtx)

type HandlerContext struct {
	container         container.Container
	encoderNegotiator NegotiatorBuilder[Encoder]
	decoderNegotiator NegotiatorBuilder[Decoder]
}

func (c *HandlerContext) RegisterEncoder(contentType string, enc Encoder, aliases ...string) {
	if c.encoderNegotiator == nil {
		c.encoderNegotiator = NewNegotiatorBuilder[Encoder]()
	}
	c.encoderNegotiator.Register(contentType, enc, aliases...)
}

func (c *HandlerContext) RegisterDecoder(contentType string, dec Decoder, aliases ...string) {
	if c.decoderNegotiator == nil {
		c.decoderNegotiator = NewNegotiatorBuilder[Decoder]()
	}
	c.decoderNegotiator.Register(contentType, dec, aliases...)
}

// H wraps your handler function with the Go generics magic.
func H[T any, O any](handle Handle[T, O], options ...Option) http.HandlerFunc {

	handlerCtx := &HandlerContext{
		container:         container.Default,
		encoderNegotiator: NewNegotiatorBuilder[Encoder](),
		decoderNegotiator: NewNegotiatorBuilder[Decoder](),
	}

	//todo: for now hardcoded, will need to make this injectable or configurable through the default
	jsonEncoder := json.NewJsonEncoding()
	handlerCtx.encoderNegotiator.Register("application/json", jsonEncoder, "*/*")
	handlerCtx.decoderNegotiator.Register("application/json", jsonEncoder)

	//xmlEncoder := xml.NewXMLEncoding()
	//handlerCtx.encoderNegotiator.Register("application/xml", xmlEncoder)
	//handlerCtx.decoderNegotiator.Register("application/xml", xmlEncoder)

	//process options
	for _, option := range options {
		option(handlerCtx)
	}

	var req T
	argumentDecoder, err := BuildArgumentsBinder(req)
	if err != nil {
		panic(err)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		enc, err := handlerCtx.encoderNegotiator.Get(req.Header.Get("Accept"))
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}

		var (
			payload = new(T)
			res     any
			e       *HTTPError
		)

		//decode arguments
		if err := argumentDecoder.Decode(req, payload); err != nil {
			e = Error(err, http.StatusBadRequest)
			goto Encode
		}

		//decode the body.
		if req.ContentLength > 0 {
			dec, err := handlerCtx.decoderNegotiator.Get(req.Header.Get("Content-Type"))
			if err != nil {
				res = err
				goto Encode
			}

			if err := dec.Decode(req, payload); err != nil {
				e = Error(err, http.StatusBadRequest)
				goto Encode
			}
		}

		//call action handler
		res, err = handle(req.Context(), *payload)
		if err != nil {
			e = Error(err, 0)
		}

	Encode:
		rw.Header().Add("Content-Type", enc.Mimetype()+"; charset=utf-8")

		if e != nil {
			res = e
		}

		if h, ok := res.(Headerer); ok {
			for k, v := range h.Header() {
				rw.Header().Add(k, v[0])
			}
		}

		if sc, ok := res.(StatusCoder); ok {
			rw.WriteHeader(sc.StatusCode())
		}

		if err = enc.Encode(rw, res); err != nil {
			log.Default().Printf("%v", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}
