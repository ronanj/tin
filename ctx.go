package tin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Context struct {
	t          *Tin
	w          http.ResponseWriter
	r          *http.Request
	params     map[string]int
	clientGone bool
}

func (t *Context) BindJSON(v interface{}) error {
	b, err := ioutil.ReadAll(t.r.Body)
	defer t.r.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (t *Context) JSON(status int, v interface{}) {

	body, _ := json.Marshal(v)
	// t.w.Header().Set("Access-Control-Allow-Origin", "*")
	t.w.Header().Set("Content-Type", "application/json")
	t.w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	t.w.Write(body)

}

func (t *Context) String(status int, v string) {

	t.w.Header().Set("Content-Length", fmt.Sprintf("%d", len(v)))
	t.w.Write([]byte(v))

}
func (t *Context) Error(err error) {
	t.JSON(http.StatusOK, H{"status": "error", "reason": err.Error()})
}

func (t *Context) Query(v string) string {
	return t.r.FormValue(v)
}

func (t *Context) GetHeader(h string) string {
	return t.r.Header.Get(h)
}

func (t *Context) Header(k, v string) {
	t.w.Header().Set(k, v)
}

func (t *Context) Param(s string) string {
	if idx, has := t.params[s]; has {
		url := strings.Split(t.r.URL.Path, "/")
		return url[idx]
	}
	panic("Invalid parameter")
}

func (t *Context) ClientIP() string {
	return t.r.RemoteAddr
}
