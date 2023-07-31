package server

import (
	"net/http"
)

type HandlerFunc func(c *Context)

var _ Server = &HTTPServer{}

type Server interface {
	Start(address string) error
}

func NewServer() *HTTPServer {
	return &HTTPServer{
		router: newRouter(),
	}
}

type HTTPServer struct {
	router
}

func (hs *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		R: request,
		W: writer,
	}
	hs.serve(ctx)
}

func (hs *HTTPServer) Post(path string, handler HandleFunc) {
	hs.addRouter(http.MethodPost, path, handler)
}

func (hs *HTTPServer) Get(path string, handler HandleFunc) {
	hs.addRouter(http.MethodGet, path, handler)
}

func (hs *HTTPServer) Start(address string) error {
	return http.ListenAndServe(address, hs)
}

func (hs *HTTPServer) serve(ctx *Context) {
	mi, ok := hs.findRoute(ctx.R.Method, ctx.R.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.W.WriteHeader(404)
		ctx.W.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = mi.pathParams
	mi.n.handler(ctx)
}
