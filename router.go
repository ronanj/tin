package tin

import (
	// "log"
	"net/http"
	"strings"
)

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

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		ctx := &Context{t, w, r, params, false}
		done := make(chan bool, 0)
		go func() {
			select {
			case <-r.Context().Done():
				ctx.clientGone = true
			case <-done:
			}
		}()

		defer func() {
			done <- true
		}()

		if r.Method == "GET" {
			handle(ctx)
		} else {
			http.Error(w, "Invalid method", http.StatusInternalServerError)
		}

	})
}

func (t *Tin) POST(path string, handle func(c *Context)) {

	path, params := extractPath(path)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		ctx := &Context{t, w, r, params, false}
		if r.Method == "POST" {
			handle(ctx)
		} else {
			http.Error(w, "Invalid method", http.StatusInternalServerError)
		}

	})
}
