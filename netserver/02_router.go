package netserver

import (
	"net/http"
	"slices"
)

type Router struct {
	globalChain []func(http.Handler) http.Handler
	routeChain  []func(http.Handler) http.Handler
	isSubRouter bool
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{ServeMux: http.NewServeMux()}
}

func (r *Router) HasRoutes() bool {
	if r.isSubRouter && len(r.routeChain) > 0 {
		return true
	} else if len(r.globalChain) > 0 {
		return true
	}

	return false
}

func (r *Router) Type() string {
	if r.isSubRouter {
		return "sub-router"
	}
	return "global-router"
}

func (r *Router) Use(mw ...func(http.Handler) http.Handler) {
	if r.isSubRouter {
		r.routeChain = append(r.routeChain, mw...)
	} else {
		r.globalChain = append(r.globalChain, mw...)
	}
}

func (r *Router) Group(fn func(r *Router)) {
	subRouter := &Router{routeChain: slices.Clone(r.routeChain), isSubRouter: true, ServeMux: r.ServeMux}
	fn(subRouter)
}

func (r *Router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.Handle(pattern, h)
}

func (r *Router) Handle(pattern string, h http.Handler) {
	for _, mw := range slices.Backward(r.routeChain) {
		h = mw(h)
	}
	r.ServeMux.Handle(pattern, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	var h http.Handler = r.ServeMux

	for _, mw := range slices.Backward(r.globalChain) {
		h = mw(h)
	}
	h.ServeHTTP(w, rq)
}

// File ends here
