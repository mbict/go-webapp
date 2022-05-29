package webapp

import (
	"net/http"
	"strings"
)

type group struct {
	r      *API
	prefix string
}

func (g *group) Get(path string, handle http.HandlerFunc) {
	g.Handle(http.MethodGet, path, handle)
}

func (g *group) Post(path string, handle http.HandlerFunc) {
	g.Handle(http.MethodPost, path, handle)
}

func (g *group) Put(path string, handle http.HandlerFunc) {
	g.Handle(http.MethodPut, path, handle)
}

func (g *group) Patch(path string, handle http.HandlerFunc) {
	g.Handle(http.MethodPatch, path, handle)
}

func (g *group) Delete(path string, handle http.HandlerFunc) {
	g.Handle(http.MethodDelete, path, handle)
}

func (g *group) Handle(method, path string, handle http.HandlerFunc) {
	g.r.Handle(method, g.prefix+path, handle)
}

func (g *group) Handler(method, path string, handle http.Handler) {
	g.r.Handler(method, g.prefix+path, handle)
}

func (g *group) Group(path string) Router {
	return &group{prefix: g.prefix + path, r: g.r}
}

func (g *group) Use(mw ...Middleware) {
	g.r.Use(func(next http.HandlerFunc) http.HandlerFunc {
		mwNext := next
		for _, fn := range mw {
			mwNext = fn(mwNext)
		}

		return func(rw http.ResponseWriter, req *http.Request) {
			// Only use the middleware if path belongs to group.
			if strings.HasPrefix(req.URL.Path, g.prefix) {
				mwNext(rw, req)
			} else {
				next(rw, req)
			}
		}
	})
}
