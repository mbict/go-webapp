package webapp

import (
	"github.com/justinas/alice"
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

func mc(mw []Middleware) []alice.Constructor {
	res := make([]alice.Constructor, len(mw))
	for i, m := range mw {
		res[i] = func(next http.Handler) http.Handler {
			return m(next.ServeHTTP)
		}
	}
	return res
}

type Router interface {
	Get(path string, handle http.HandlerFunc, mw ...Middleware)
	Post(path string, handle http.HandlerFunc, mw ...Middleware)
	Put(path string, handle http.HandlerFunc, mw ...Middleware)
	Patch(path string, handle http.HandlerFunc, mw ...Middleware)
	Delete(path string, handle http.HandlerFunc, mw ...Middleware)
	Handle(method, path string, handle http.HandlerFunc, mw ...Middleware)
	Handler(method, path string, handle http.Handler, mw ...Middleware)
	Group(path string, mw ...Middleware) Router

	//Use(mw ...Middleware)
}

type API struct {
	router     *httprouter.Router
	config     *Config
	container  container.Container
	middleware alice.Chain

	NotFound         http.HandlerFunc
	MethodNotAllowed http.HandlerFunc
	PanicHandler     func(http.ResponseWriter, *http.Request, interface{})
}

func (r *API) Get(path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handle(http.MethodGet, path, handle, mw...)
}

func (r *API) Post(path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handle(http.MethodPost, path, handle, mw...)
}

func (r *API) Put(path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handle(http.MethodPut, path, handle, mw...)
}

func (r *API) Patch(path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handle(http.MethodPatch, path, handle, mw...)
}

func (r *API) Delete(path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handle(http.MethodDelete, path, handle, mw...)
}

func (r *API) Handle(method, path string, handle http.HandlerFunc, mw ...Middleware) {
	r.Handler(method, path, handle, mw...)
}

func (r *API) Handler(method, path string, handle http.Handler, mw ...Middleware) {
	r.router.Handler(method, path, alice.New(mc(mw)...).Then(handle))
}

func (r *API) Group(path string, mw ...Middleware) Router {
	return &group{
		prefix:     path,
		r:          r,
		middleware: alice.New(mc(mw)...),
	}
}

// global middleware
func (r *API) Use(mw ...Middleware) {
	for _, m := range mw {
		r.middleware = r.middleware.Append(func(handler http.Handler) http.Handler {
			return m(handler.ServeHTTP)
		})
	}
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
		middleware:       alice.New(),
	}
}

func (r *API) RequestHander() http.HandlerFunc {
	r.router.NotFound = r.NotFound
	r.router.MethodNotAllowed = r.MethodNotAllowed
	//r.PanicHandler = r.PanicHandler

	h := r.middleware.Then(r.router).ServeHTTP

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
