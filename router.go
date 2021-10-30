package tin

import (
	"fmt"
	"net/http"
)

type route struct {
	path    *path
	handler http.Handler
	method  string
}

type tinRouter struct {
	routes []*route
}

func newTinRouter() *tinRouter {
	return &tinRouter{
		routes: make([]*route, 0),
	}
}

func (h *tinRouter) add(path *path, method string, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{path, http.HandlerFunc(handler), method})
}

func (h *tinRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if route, err := h.findRoute(r.URL.Path, r.Method); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if route == nil {
		http.NotFound(w, r)
	} else {
		route.handler.ServeHTTP(w, r)
	}

}

/* -------------------------------- */

func (h *tinRouter) findRoute(url string, method string) (*route, error) {

	invalidMethod := false
	for _, route := range h.routes {
		if route.path.match(url) {
			if method == route.method || route.method == "" {
				return route, nil
			} else {
				invalidMethod = true
			}
		}
	}

	if invalidMethod {
		return nil, fmt.Errorf("invalid method")
	}
	return nil, nil
}

func (t *Tin) GET(url string, handle func(c *Context)) {

	path := extractPath(url)

	t.router.add(path, "GET", func(w http.ResponseWriter, r *http.Request) {

		ctx := t.newContext(w, r, path)
		done := make(chan bool)
		go func() {
			select {
			case <-r.Context().Done():
				ctx.clientGone = true
			case <-done:
			}
		}()

		defer func() {
			if !ctx.clientGone {
				done <- true
			}
		}()

		t.handle(handle, ctx)

	})
}

func (t *Tin) POST(url string, handle func(c *Context)) {

	path := extractPath(url)

	t.router.add(path, "POST", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, path))

	})
}

func (t *Tin) DELETE(url string, handle func(c *Context)) {

	path := extractPath(url)

	t.router.add(path, "DELETE", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, path))

	})
}

func (t *Tin) Any(url string, handle func(c *Context)) {

	path := extractPath(url)

	t.router.add(path, "", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, path))

	})
}
