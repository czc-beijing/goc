package server

import (
	"fmt"
	"goc/context"
	"net/http"
)

type Server interface {
	Routable
	Start(address string) error
}

type HttpServer struct {
	router Router
}

func New() *HttpServer {
	return &HttpServer{
		router: NewRouter(),
	}
}

func (hs *HttpServer) Router(method string, pattern string, handler func(c *context.Context)) {
	hs.router.Router(method, pattern, handler)
}

func (hs *HttpServer) Start(address string) error {
	fmt.Println(address, "run........")
	return http.ListenAndServe(address, hs.router)
}
