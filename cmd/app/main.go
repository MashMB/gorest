package main

import (
	"github.com/mashmb/gorest/internal/app/settings"
	"github.com/mashmb/gorest/internal/app/web"
)

func main() {
	settings := settings.LoadSettings()
	server := web.NewServer(settings)
	server.Run()
}
