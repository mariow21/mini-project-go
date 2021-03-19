package http

import (
	"net/http"

	"testing/pkg/grace"

	"github.com/rs/cors"
)

// TestingHandler ...
type TestingHandler interface {
	// Masukkan fungsi handler di sini
	TestingHandler(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	server  *http.Server
	Testing TestingHandler
}

// Serve is serving HTTP gracefully on port x ...
func (s *Server) Serve(port string) error {
	handler := cors.AllowAll().Handler(s.Handler())
	return grace.Serve(port, handler)
}
