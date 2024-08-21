package main

import "github.com/mashmb/gorest/internal/app/web"

func main() {
	server := web.NewServer()
	server.Run()
}
