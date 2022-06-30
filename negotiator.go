package webapp

import "strings"

type Negotiator[T any] interface {
	Get(mimetype string) (T, error)
}

type NegotiatorBuilder[T any] interface {
	Negotiator[T]

	Register(mimetype string, v T, aliases ...string)
}

type negotiator[T any] struct {
	encodings map[string]T
	aliases   map[string]string
}

func stripEncoding(s string) string {
	index := strings.IndexRune(s, ';')
	if index > 0 {
		s = s[:index]
	}
	return strings.TrimSpace(s)
}

func (n negotiator[T]) Get(contentTypes string) (h T, e error) {
	for _, mimetype := range strings.Split(contentTypes, ",") {
		mimetype = stripEncoding(mimetype)

		if v, ok := n.encodings[mimetype]; ok {
			return v, nil
		}

		if name := n.aliases[mimetype]; name != "" {
			return n.encodings[name], nil
		}
	}
	return h, ErrNotAcceptable
}

func (n negotiator[T]) Register(mimetype string, v T, aliases ...string) {
	n.encodings[mimetype] = v
	for _, alias := range aliases {
		n.aliases[alias] = mimetype
	}
}

func NewNegotiator[T any]() Negotiator[T] {
	return NewNegotiatorBuilder[T]()
}

func NewNegotiatorBuilder[T any]() NegotiatorBuilder[T] {
	return &negotiator[T]{
		encodings: make(map[string]T),
		aliases:   make(map[string]string),
	}
}
