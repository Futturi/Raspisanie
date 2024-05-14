package service

import (
	"net/http"
	"time"
)

type Server struct {
	Server *http.Server
}

func (s *Server) InitServer(port string, handler http.Handler) error {
	s.Server = &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s.Server.ListenAndServe()
}
