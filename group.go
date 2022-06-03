package webapp

import (
	"github.com/justinas/alice"
	"net/http"
)

type group struct {
	r          *API
	middleware alice.Chain
	prefix     string
}

func (g *group) Get(path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(http.MethodGet, path, handle, mw...)
}

func (g *group) Post(path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(http.MethodPost, path, handle, mw...)
}

func (g *group) Put(path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(http.MethodPut, path, handle, mw...)
}

func (g *group) Patch(path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(http.MethodPatch, path, handle, mw...)
}

func (g *group) Delete(path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(http.MethodDelete, path, handle, mw...)
}

func (g *group) Handle(method, path string, handle http.HandlerFunc, mw ...Middleware) {
	g.Handler(method, path, handle, mw...)
}

func (g *group) Handler(method, path string, handle http.Handler, mw ...Middleware) {
	g.r.Handler(method, g.prefix+path, g.middleware.Append(mc(mw)...).Then(handle))
}

func (g *group) Group(path string, mw ...Middleware) Router {
	return &group{prefix: g.prefix + path, r: g.r, middleware: g.middleware.Append(mc(mw)...)}
}

//
//func (g *group) Use(mw ...Middleware) {
//	g.r.Use(func(next http.HandlerFunc) http.HandlerFunc {
//		mwNext := next
//		for _, fn := range mw {
//			mwNext = fn(mwNext)
//		}
//
//		return func(rw http.ResponseWriter, req *http.Request) {
//			// Only use the middleware if path belongs to group.
//			if strings.HasPrefix(req.URL.Path, g.prefix) {
//				mwNext(rw, req)
//			} else {
//				next(rw, req)
//			}
//		}
//	})
//}
