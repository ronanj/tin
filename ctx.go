package tin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Context struct {
	t       *Tin
	Writer  http.ResponseWriter
	Request *http.Request
	path    *path

	clientGone bool
	isAborted  bool /* Used by the middleware */
}

func (t *Tin) newContext(w http.ResponseWriter, r *http.Request, path *path) *Context {

	return &Context{
		t:          t,
		Writer:     w,
		Request:    r,
		path:       path,
		clientGone: false,
	}
}

func (t *Context) BindJSON(v interface{}) error {
	b, err := ioutil.ReadAll(t.Request.Body)
	defer t.Request.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func (t *Context) JSON(status int, v interface{}) {

	body, _ := json.Marshal(v)
	t.Writer.Header().Set("Content-Type", "application/json")
	t.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	t.Writer.WriteHeader(status)
	t.Writer.Write(body)

}

func (t *Context) String(status int, v string) {

	t.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(v)))
	t.Writer.Write([]byte(v))

}
func (t *Context) Error(err error) {
	t.JSON(http.StatusOK, H{"status": "error", "reason": err.Error()})
}

func (t *Context) Query(v string) string {
	return t.Request.FormValue(v)
}

func (t *Context) GetHeader(h string) string {
	return t.Request.Header.Get(h)
}

func (t *Context) Header(k, v string) {
	t.Writer.Header().Set(k, v)
}

func (t *Context) Param(s string) string {
	return t.path.param(t.Request.URL.Path, s)
}

func (t *Context) ClientIP() string {
	return t.Request.RemoteAddr
}
