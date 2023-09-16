package tin

import (
	"net/http"
)

type route struct {
	path    *path
	handler func(ctx *Context)
	method  string
}

type tinRouter struct {
	tin    *Tin
	routes []*route
}

func newTinRouter(t *Tin) *tinRouter {
	return &tinRouter{t, make([]*route, 0)}
}

func (h *tinRouter) add(path *path, method string, handler func(ctx *Context)) {
	h.routes = append(h.routes, &route{path, handler, method})
}

func (h *tinRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	route, validMethod, err := h.findRoute(r.URL.Path, r.Method)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var path *path
	if route != nil {
		path = route.path
	}

	ctx := newContext(w, r, path)
	if !h.tin.applyMiddleware(ctx) {
		return
	}

	if ctx.activateRecovery {
		defer func() {
			if err := recover(); err != nil {
				if ctx.recoveryNotifier != nil {
					ctx.recovery(err, !ctx.recoveryNotifier(ctx, err))

				} else {
					ctx.recovery(err, true)
				}
			}
		}()
	}

	if route == nil {
		http.NotFound(w, r)
	} else if validMethod {
		route.handler(ctx)
	} else {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
	}

}

/* -------------------------------- */

func (h *tinRouter) findRoute(url string, method string) (*route, bool, error) {

	for _, route := range h.routes {
		if route.path.match(url) {
			validMethod := method == route.method || route.method == ""
			return route, validMethod, nil
		}
	}

	return nil, false, nil
}

func (t *Tin) GET(url string, handler func(c *Context)) {

	path := extractPath(url)

	t.router.add(path, "GET", func(ctx *Context) {

		done := make(chan bool)
		go func() {
			select {
			case <-ctx.Request.Context().Done():
				ctx.clientGone = true
			case <-done:
			}
		}()

		defer func() {
			if !ctx.clientGone {
				done <- true
			}
		}()

		handler(ctx)

	})
}

func (t *Tin) POST(url string, handler func(c *Context)) {

	path := extractPath(url)
	t.router.add(path, "POST", handler)
}

func (t *Tin) DELETE(url string, handler func(c *Context)) {

	path := extractPath(url)
	t.router.add(path, "DELETE", handler)

}

func (t *Tin) Any(url string, handler func(c *Context)) {

	path := extractPath(url)
	t.router.add(path, "", handler)

}
