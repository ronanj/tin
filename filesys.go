package tin

import (
	"net/http"
	"strings"
)

func (t *Tin) Static(relativePath, root string) {
	t.StaticFS(relativePath, http.Dir(root))
}

func (t *Tin) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := t.createStaticHandler(relativePath, fs)
	t.GET(relativePath, handler)
}

func (t *Tin) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {

	fileServer := http.StripPrefix(relativePath, http.FileServer(fs))

	return func(c *Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
