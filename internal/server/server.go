package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(bind Host, port Port) error {
	addr := fmt.Sprintf("%s:%d", bind, port)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Printf("[HTTP] Listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[HTTP] Failed to start server: %v", err)
	}

	return nil
}
