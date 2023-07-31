package main

import (
	server2 "goc/server"
)

func main() {
	server := server2.NewServer()
	server.Get("/hello", server2.Hello)
	_ = server.Start(":8090")
}
