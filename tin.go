package tin

import (
	"net/http"
)

type Tin struct {
}

type H = map[string]interface{}

func Default() *Tin {
	return &Tin{}
}

func (t *Tin) Run(address string) {

	http.ListenAndServe(address, nil)

}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func SetMode(_ string) {
}
