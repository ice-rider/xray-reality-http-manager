package http

import (
	"fmt"
	"net/http"
	"xray_server/internal/app"
)

type Server struct {
	configService *app.ConfigService
	port          int
}

func NewServer(configService *app.ConfigService, port int) *Server {
	return &Server{
		configService: configService,
		port:          port,
	}
}

func (s *Server) Start() error {
	handler := NewClientHandler(s.configService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/client", handler.CreateClient)

	addr := fmt.Sprintf(":%d", s.port)

	return http.ListenAndServe(addr, mux)
}
