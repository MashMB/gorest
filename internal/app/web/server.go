package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mashmb/gorest/internal/app/settings"
)

type server struct {
	settings settings.Settings
}

func NewServer(stg settings.Settings) *server {
	return &server{
		settings: stg,
	}
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", s.settings.Server.Host, s.settings.Server.Port)
	server := &http.Server{
		Addr: addr,
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
