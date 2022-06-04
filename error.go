package webapp

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/mbict/go-webapp/encoding/json"
	"log"
	"net/http"
	"strconv"
)

func E(err error) http.HandlerFunc {
	h := H(func(context.Context, Empty) (*Empty, error) {
		return nil, err
	})
	return func(rw http.ResponseWriter, req *http.Request) {
		h(rw, req)
	}
}

func ErrorHandler() func(error) http.HandlerFunc {
	encoderNegotiator := NewNegotiatorBuilder[Encoder]()
	jsonEncoder := json.NewJsonEncoding()
	encoderNegotiator.Register("application/json", jsonEncoder, "*/*")

	return func(e error) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			enc, err := encoderNegotiator.Get(req.Header.Get("Accept"))
			//default to json encoder if no negotiation cna be done
			if err != nil {
				enc = jsonEncoder
			}

			rw.Header().Add("Content-Type", enc.Mimetype()+"; charset=utf-8")

			if h, ok := e.(Headerer); ok {
				for k, v := range h.Header() {
					rw.Header().Add(k, v[0])
				}
			}

			if sc, ok := e.(StatusCoder); ok {
				rw.WriteHeader(sc.StatusCode())
			}

			if err = enc.Encode(rw, e); err != nil {
				log.Default().Printf("%v", err)
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}
}

type HTTPError struct {
	err  error
	code int
}

func Error(err error, code int) error {
	if _, ok := err.(StatusCoder); ok {
		return err
	}
	return &HTTPError{err, code}
}

func (e *HTTPError) Error() string {
	return e.err.Error()
}

func (e *HTTPError) StatusCode() int {
	if sc, ok := e.err.(StatusCoder); ok {
		return sc.StatusCode()
	}

	if e.code == 0 {
		return http.StatusInternalServerError
	}

	return e.code
}

func (e *HTTPError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

func (e *HTTPError) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"message":`)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e *HTTPError) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}

func (e *HTTPError) MarshalYAML() (interface{}, error) {
	return map[string]string{"message": e.Error()}, nil
}

//----

// StatusError creates an error from an HTTP status code.
type StatusError int

const (
	ErrBadRequest           = StatusError(http.StatusBadRequest)
	ErrUnauthorized         = StatusError(http.StatusUnauthorized)
	ErrForbidden            = StatusError(http.StatusForbidden)
	ErrNotFound             = StatusError(http.StatusNotFound)
	ErrMethodNotAllowed     = StatusError(http.StatusMethodNotAllowed)
	ErrNotAcceptable        = StatusError(http.StatusNotAcceptable)
	ErrUnsupportedMediaType = StatusError(http.StatusUnsupportedMediaType)
	ErrInternalServerError  = StatusError(http.StatusInternalServerError)
)

func (e StatusError) Error() string {
	return http.StatusText(int(e))
}

func (e StatusError) StatusCode() int {
	return int(e)
}

func (e StatusError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

func (e StatusError) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString(`{"message":`)
	buf.WriteString(strconv.Quote(e.Error()))
	buf.WriteRune('}')

	return buf.Bytes(), nil
}

func (e StatusError) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: "message"}}
	return enc.EncodeElement(e.Error(), start)
}

func (e StatusError) MarshalYAML() (interface{}, error) {
	return map[string]string{"message": e.Error()}, nil
}
