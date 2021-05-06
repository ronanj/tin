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
	return &Tin{
		router:      newTinRouter(),
		middlewares: make([]HandlerFunc, 0),
	}
}

func New() *Tin {
	return &Tin{
		router:      newTinRouter(),
		middlewares: make([]HandlerFunc, 0),
	}
}

func (t *Tin) Use(middleware HandlerFunc) {
	t.middlewares = append(t.middlewares, middleware)
}

func (t *Tin) Run(address string) {

	http.ListenAndServe(address, t.router)

}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func SetMode(_ string) {
}

type HandlerFunc = func(c *Context)
