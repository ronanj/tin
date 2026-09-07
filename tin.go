package tin

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type Tin struct {
	router      *tinRouter
	middlewares []HandlerFunc
	server      http.Server
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

	t.server = http.Server{
		Addr:    address,
		Handler: t.router,
	}

	var err error
	var listener net.Listener
	if strings.HasPrefix(address, "unix:") {
		listener, err = net.Listen("unix", address[5:])
	} else {
		listener, err = net.Listen("tcp", address)
	}
	if err != nil {
		return err
	}

	return t.server.Serve(listener)

}

func (t *Tin) Shutdown(ctx context.Context) error {
	return t.server.Shutdown(ctx)
}

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

func SetMode(_ string) {
}

type HandlerFunc = func(c *Context)
