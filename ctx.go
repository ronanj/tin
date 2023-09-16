package tin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	path    *path

	headerWritten    bool
	clientGone       bool
	isAborted        bool /* Used by the middleware */
	activateRecovery bool
	recoveryNotifier func(*Context, interface{}) bool
}

func newContext(w http.ResponseWriter, r *http.Request, path *path) *Context {

	return &Context{
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
	t.headerWritten = true

}

func (t *Context) String(status int, v string) {

	t.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(v)))
	t.Writer.Write([]byte(v))

}
func (t *Context) Error(err error) {
	t.JSON(http.StatusOK, H{"status": "error", "reason": err.Error()})
}

func (c *Context) Abort() {
	c.isAborted = true
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

func (t *Context) PostForm(s string) string {
	return t.Request.FormValue(s)
}

func (c *Context) AbortWithError(code int, err error) {
	c.AbortWithStatus(code)
	c.Error(err)
	c.headerWritten = true
}

// DefaultQuery returns the keyed url query value if it exists,
// otherwise it returns the specified defaultValue string.
// See: Query() and GetQuery() for further information.
//
//	GET /?name=Manu&lastname=
//	c.DefaultQuery("name", "unknown") == "Manu"
//	c.DefaultQuery("id", "none") == "none"
//	c.DefaultQuery("lastname", "none") == ""
func (c *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := c.Request.URL.Query()[key]; ok {
		return value[0]
	}
	return defaultValue
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) {
	http.Redirect(c.Writer, c.Request, location, code)
	c.headerWritten = true
}

// Data writes some data into the body stream and updates the HTTP code.
func (c *Context) Data(code int, contentType string, data []byte) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set("Content-Type", contentType)
	c.Writer.Write(data)
	c.headerWritten = true
}
