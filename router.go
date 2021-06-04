package tin

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type route struct {
	pattern *regexp.Regexp
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

func (h *tinRouter) add(path string, method string, handler func(http.ResponseWriter, *http.Request)) {
	pattern, err := regexp.Compile(path)
	if err != nil {
		log.Fatal("Invalid path", path, err)
	}
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler), method})
}

func (h *tinRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	invalidMethod := false
	for _, route := range h.routes {

		if route.pattern.MatchString(r.URL.Path) {
			if r.Method == route.method || route.method == "" {

				route.handler.ServeHTTP(w, r)
				return

			} else {
				invalidMethod = true
			}
		}
	}

	if invalidMethod {
		http.Error(w, "Invalid method", http.StatusInternalServerError)
	}

	http.NotFound(w, r)
}

/* -------------------------------- */

func extractPath(path string) (string, map[string]int) {

	if !strings.Contains(path, "/:") {
		return path, nil
	}

	parts := make([]string, 0)
	params := make(map[string]int)
	for i, part := range strings.Split(path, "/") {

		// log.Println(path, ">", i, parts)
		if len(part) > 0 && part[0] == ':' {
			params[part[1:]] = i
		} else {
			if len(params) > 0 {
				panic("Invalid path")
			}
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, "/") + "/", params

}

func (t *Tin) GET(path string, handle func(c *Context)) {

	path, params := extractPath(path)

	t.router.add(path, "GET", func(w http.ResponseWriter, r *http.Request) {

		ctx := t.newContext(w, r, params)
		done := make(chan bool, 0)
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

func (t *Tin) POST(path string, handle func(c *Context)) {

	path, params := extractPath(path)

	t.router.add(path, "POST", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, params))

	})
}


func (t *Tin) DELETE(path string, handle func(c *Context)) {

	path, params := extractPath(path)

	t.router.add(path, "DELETE", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, params))

	})
}

func (t *Tin) Any(path string, handle func(c *Context)) {

	path, params := extractPath(path)

	t.router.add(path, "", func(w http.ResponseWriter, r *http.Request) {

		t.handle(handle, t.newContext(w, r, params))

	})
}
