package tin

import (
	"net/http"
)

type Tin struct {
	router *tinRouter
}

type H = map[string]interface{}

func Default() *Tin {
	return &Tin{}
}

func New() *Tin {
	return &Tin{router: newTinRouter()}
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
