package tin

import (
	"net/http"
)

type Tin struct {
	router      *tinRouter
	middlewares []HandlerFunc
}

type H = map[string]interface{}

func Default() *Tin {
	t := &Tin{}
	t.router = newTinRouter(t)
	t.middlewares = make([]HandlerFunc, 0)
	return t
}

func New() *Tin {
	return Default()
}

func (t *Tin) Run(address string) error {

	return http.ListenAndServe(address, t.router)

}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func SetMode(_ string) {
}

type HandlerFunc = func(c *Context)
