package http

import (
	"fmt"
	"net/http"
)

type Server struct {
	authHandler    *AuthHandler
	clientHandler  *ClientHandler
	statsHandler   *StatsHandler
	authMiddleware *AuthMiddleware
	port           int
}

func NewServer(
	authHandler *AuthHandler,
	clientHandler *ClientHandler,
	statsHandler *StatsHandler,
	authMiddleware *AuthMiddleware,
	port int,
) *Server {
	return &Server{
		authHandler:    authHandler,
		clientHandler:  clientHandler,
		statsHandler:   statsHandler,
		authMiddleware: authMiddleware,
		port:           port,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/login", s.authHandler.Login)

	mux.Handle("/api/v1/client", s.authMiddleware.Middleware(http.HandlerFunc(s.clientHandler.CreateClient)))
	mux.Handle("/api/v1/clients", s.authMiddleware.Middleware(http.HandlerFunc(s.clientHandler.GetClients)))
	mux.Handle("/api/v1/stats", s.authMiddleware.Middleware(http.HandlerFunc(s.statsHandler.GetClientsStats)))
	mux.Handle("/api/v1/stats/email", s.authMiddleware.Middleware(http.HandlerFunc(s.statsHandler.GetClientStats)))

	addr := fmt.Sprintf(":%d", s.port)

	return http.ListenAndServe(addr, mux)
}
