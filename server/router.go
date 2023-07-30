package router

import (
	"fmt"
	"goc/context"
	"net/http"
)

type Routable interface {
	Router(method string, pattern string, handler func(c *context.Context))
}

type Router interface {
	http.Handler
	Routable
}

type HttpRouter struct {
	trees map[string]func(c *context.Context)
}

func New() *HttpRouter {
	return &HttpRouter{}
}

func (r *HttpRouter) Route(method string, pattern string, handler func(c *context.Context)) {
	r.trees[r.key(method, pattern)] = handler
}

func (r *HttpRouter) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := r.key(request.Method, request.URL.Path)
	if handler, ok := r.trees[key]; ok {
		handler(context.NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("no method fond"))
	}
}

func (r *HttpRouter) key(method string, pattern string) string {
	return fmt.Sprintf("%s#%s", method, pattern)
}
