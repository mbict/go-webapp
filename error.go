package webapp

import (
	"bytes"
	"context"
	"encoding/xml"
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

type HTTPError struct {
	err  error
	code int
}

func Error(err error, code int) *HTTPError {
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
