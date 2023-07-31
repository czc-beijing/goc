package main

import (
	"goc/web"
)

func main() {
	server := web.NewHTTPServer()
	server.Get("/hello", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})
	_ = server.Start(":8090")
}
