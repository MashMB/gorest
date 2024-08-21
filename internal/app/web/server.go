package web

import (
	"fmt"
	"log/slog"
	"net/http"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	addr := fmt.Sprintf("%s:%s", "0.0.0.0", "8080")
	server := &http.Server{
		Addr: addr,
	}
	slog.Info("HTTP server started", "addr", addr)
	server.ListenAndServe()
}
