package main

import (
	"goc/api"
	server2 "goc/server"
	"net/http"
)

func main() {
	server := server2.New()
	server.Router(http.MethodGet, "/hello", api.Hello)
	_ = server.Start(":8090")
}
