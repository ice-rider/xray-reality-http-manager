package http

import (
	"fmt"
	"net/http"
)

type Server struct {
	authHandler      *AuthHandler
	clientHandler    *ClientHandler
	statsHandler     *StatsHandler
	mldsa65Handler   *Mldsa65Handler
	configHandler    *ServerConfigHandler
	authMiddleware   *AuthMiddleware
	port             int
}

func NewServer(
	authHandler *AuthHandler,
	clientHandler *ClientHandler,
	statsHandler *StatsHandler,
	mldsa65Handler *Mldsa65Handler,
	configHandler *ServerConfigHandler,
	authMiddleware *AuthMiddleware,
	port int,
) *Server {
	return &Server{
		authHandler:      authHandler,
		clientHandler:    clientHandler,
		statsHandler:     statsHandler,
		mldsa65Handler:   mldsa65Handler,
		configHandler:    configHandler,
		authMiddleware:   authMiddleware,
		port:             port,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/login", s.authHandler.Login)

	mux.Handle("/api/v1/client", s.authMiddleware.Middleware(http.HandlerFunc(s.clientHandler.CreateClient)))
	mux.Handle("/api/v1/clients", s.authMiddleware.Middleware(http.HandlerFunc(s.clientHandler.GetClients)))
	mux.Handle("/api/v1/stats", s.authMiddleware.Middleware(http.HandlerFunc(s.statsHandler.GetClientsStats)))
	mux.Handle("/api/v1/stats/email", s.authMiddleware.Middleware(http.HandlerFunc(s.statsHandler.GetClientStats)))
	mux.Handle("/api/v1/mldsa65-public", s.authMiddleware.Middleware(http.HandlerFunc(s.mldsa65Handler.GetMldsa65Public)))
	mux.Handle("/api/v1/server-config", s.authMiddleware.Middleware(http.HandlerFunc(s.configHandler.GetServerConfig)))

	addr := fmt.Sprintf(":%d", s.port)

	return http.ListenAndServe(addr, mux)
}
