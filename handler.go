package webapp

//https://github.com/abemedia/go-don

import (
	"context"
	"github.com/mbict/go-webapp/container"
	"github.com/mbict/go-webapp/decoder"
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
	defaultEncoding   string
	errorHandler      func(error) error
}

func (ctx *HandlerContext) getEncoder(acceptMimetype string) (Encoder, error) {
	enc, err := ctx.encoderNegotiator.Get(acceptMimetype)
	if err != nil {
		//try to get the default encoder
		if enc, err = ctx.encoderNegotiator.Get(ctx.defaultEncoding); err != nil {
			//if all fail we use the defaultEncoder
			enc = defaultEncoder
			err = nil
		}
	}
	return enc, err
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

/* todo swap out for text encoding, this is used to output errors if accept type is notavailable */
var defaultEncoder = json.NewJsonEncoding()

// H wraps your handler function with the Go generics magic.
func H[T any, O any](handle Handle[T, O], options ...Option) http.HandlerFunc {

	//create the default configurable context
	handlerCtx := &HandlerContext{
		container:         container.Default,
		encoderNegotiator: NewNegotiatorBuilder[Encoder](),
		decoderNegotiator: NewNegotiatorBuilder[Decoder](),
		defaultEncoding:   "application/json",
		errorHandler: func(err error) error {
			if _, ok := err.(StatusCoder); ok {
				return err
			}
			return Error(err, http.StatusInternalServerError)
		},
	}

	//if no options are provided use the global default ones
	if len(options) == 0 {
		options = DefaultOptions
	}

	//process options
	for _, option := range options {
		option(handlerCtx)
	}

	//internal handler for rendering errors
	handleError := func(e error, rw http.ResponseWriter, req *http.Request) {
		e = handlerCtx.errorHandler(e)

		enc, err := handlerCtx.getEncoder(req.Header.Get("Accept"))
		if err != nil {
			//cannot recover from this one
			panic("cannot determine the request encoder")
		}

		rw.Header().Add("Content-Type", enc.Mimetype()+"; charset=utf-8")

		if h, ok := e.(Headerer); ok {
			for k, v := range h.Header() {
				rw.Header().Add(k, v[0])
			}
		}

		if sc, ok := e.(StatusCoder); ok {
			rw.WriteHeader(sc.StatusCode())
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		if err = enc.Encode(rw, e); err != nil {
			log.Printf("unable to encode error in error handler %v", err)
		}
	}

	var req T

	argumentDecoder, err := BuildArgumentsBinder(req)
	if err != nil {
		panic(err)
	}

	var defaultsDecoder decoder.Decode
	if hasTag(req, defaultTag) {
		defaultsDecoder, err = decoder.NewDefaultDecoder(req, defaultTag)
		if err != nil {
			panic(err)
		}
	}

	isEmpty := makeEmptyCheck(*new(O))

	return func(rw http.ResponseWriter, req *http.Request) {
		enc, err := handlerCtx.getEncoder(req.Header.Get("Accept"))
		if err != nil {
			handleError(err, rw, req)
			return
		}

		var (
			payload = new(T)
			res     any
		)

		//set defaults
		if defaultsDecoder != nil {
			if err := defaultsDecoder(req, payload); err != nil {
				handleError(Error(err, http.StatusBadRequest), rw, req)
				return
			}
		}

		//decode the body.
		if req.ContentLength > 0 {
			dec, err := handlerCtx.decoderNegotiator.Get(req.Header.Get("Content-Type"))
			if err != nil {
				handleError(err, rw, req)
				return
			}

			if err := dec.Decode(req, payload); err != nil {
				handleError(Error(err, http.StatusBadRequest), rw, req)
				return
			}
		}

		//decode arguments, last as this should always override previous body decode
		if err := argumentDecoder.Decode(req, payload); err != nil {
			handleError(Error(err, http.StatusBadRequest), rw, req)
			return
		}

		//call action handler
		res, err = handle(req.Context(), *payload)
		if err != nil {
			handleError(err, rw, req)
			return
		}

		rw.Header().Add("Content-Type", enc.Mimetype()+"; charset=utf-8")

		if h, ok := res.(Headerer); ok {
			for k, v := range h.Header() {
				rw.Header().Add(k, v[0])
			}
		}

		if sc, ok := res.(StatusCoder); ok {
			rw.WriteHeader(sc.StatusCode())
		}

		if false == isEmpty(res) {
			if err = enc.Encode(rw, res); err != nil {
				handleError(Error(err, http.StatusInternalServerError), rw, req)
			}
		}
	}
}
