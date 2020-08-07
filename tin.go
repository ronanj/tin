package tin

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Tin struct {
}

type tinH = map[string]interface{}

func tinNew() *Tin {
	return &Tin{}

}

func (t *Tin) Run(address string) {

	http.ListenAndServe(address, nil)

}

type tinContext struct {
	t          *Tin
	w          http.ResponseWriter
	r          *http.Request
	clientGone bool
}

func (t *Tin) GET(path string, handle func(c *tinContext)) {

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {

		ctx := &tinContext{t, w, r, false}
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

func (t *tinContext) JSON(status int, v interface{}) {

	body, _ := json.Marshal(v)
	t.w.Header().Set("Access-Control-Allow-Origin", "*")
	t.w.Header().Set("Content-Type", "application/json")
	t.w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	t.w.Write(body)

}

func (t *tinContext) Error(err error) {
	t.JSON(http.StatusOK, tinH{"status": "error", "reason": err.Error()})
}

func (t *tinContext) Query(v string) string {
	return t.r.FormValue(v)
}
