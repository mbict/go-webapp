package webapp

import (
	"github.com/mbict/go-webapp/container"
	"github.com/mbict/httprouter"
	"net/http"
	"strings"
)

var DefaultEncoding = "application/json"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func WrapMiddleware(handler func(next http.Handler) http.Handler) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return handler(next).ServeHTTP
	}
}

type Router interface {
	Get(path string, handle http.HandlerFunc)
	Post(path string, handle http.HandlerFunc)
	Put(path string, handle http.HandlerFunc)
	Patch(path string, handle http.HandlerFunc)
	Delete(path string, handle http.HandlerFunc)
	Handle(method, path string, handle http.HandlerFunc)
	Handler(method, path string, handle http.Handler)
	Group(path string) Router
	Use(mw ...Middleware)
}

type API struct {
	router    *httprouter.Router
	config    *Config
	container container.Container
	mw        []Middleware

	NotFound         http.HandlerFunc
	MethodNotAllowed http.HandlerFunc
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func (r *API) Get(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *API) Post(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPost, path, handle)
}

func (r *API) Put(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPut, path, handle)
}

func (r *API) Patch(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodPatch, path, handle)
}

func (r *API) Delete(path string, handle http.HandlerFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *API) Handle(method, path string, handle http.HandlerFunc) {
	r.router.Handler(method, path, handle)
}

func (r *API) Handler(method, path string, handle http.Handler) {
	r.router.Handler(method, path, handle)
}

func (r *API) Group(path string) Router {
	return &group{prefix: path, r: r}
}

func (r *API) Use(mw ...Middleware) {
	r.mw = append(r.mw, mw...)
}

type Config struct {
	DefaultEncoding string

	//// DisableNoContent controls whether a nil or zero value response should
	//// automatically return 204 No Content with an empty body.
	//DisableNoContent bool
}

// New creates a new API instance.
func New(c container.Container) *API {
	//if c == nil {
	//	c = c
	//}

	//if c.DefaultEncoding == "" {
	//	c.DefaultEncoding = DefaultEncoding
	//}

	r := httprouter.New()

	return &API{
		router: r,
		config: &Config{
			DefaultEncoding: DefaultEncoding,
		},
		container:        c,
		NotFound:         http.NotFound,
		MethodNotAllowed: E(ErrMethodNotAllowed),
	}
}

func (r *API) RequestHander() http.HandlerFunc {
	r.router.NotFound = r.NotFound
	r.router.MethodNotAllowed = r.MethodNotAllowed
	r.PanicHandler = r.PanicHandler

	h := r.router.ServeHTTP
	for _, mw := range r.mw {
		h = mw(h)
	}

	return func(rw http.ResponseWriter, req *http.Request) {
		//set default encoding, for content type if none is set
		ct := req.Header.Get("Content-Type")
		if len(ct) == 0 || strings.HasPrefix(ct, "*/*") {
			req.Header.Set("Content-Type", r.config.DefaultEncoding)
		}

		//set default encoding for accept if none is set
		ac := req.Header.Get("Accept")
		if len(ac) == 0 || strings.HasPrefix(ac, "*/*") {
			req.Header.Set("Accept", r.config.DefaultEncoding)
		}

		h(rw, req)
	}
}

func (r *API) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, r.RequestHander())
}
