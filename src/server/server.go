package server

import (
	"net/http"
)

type Server struct {
	serv *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		serv: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.serv.ListenAndServe()
}
