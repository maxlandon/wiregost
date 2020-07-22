package https

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server - Handles all HTTP(S) ghost implant communications.
// This object has methods defined in other files of his package, for instance
// implant session management is defined in 'session.go'.
// Functions in this file are pretty much restricted to server setup and registration.
type Server struct{}

// GetServerBanner - Returns the banner for server pings
func (s *Server) GetServerBanner() (banner string) {
	return
}

// GetPoweredByBanner - Returns the banner for server pings
func (s *Server) GetPoweredByBanner() (banner string) {
	return
}

// SetupRouter - Prepares and register all HTTP handler functions for implant communications
func (s *Server) SetupRouter() (router *mux.Router) {
	return
}

// DefaultRespHeaders - Configures default HTTP response headers
func (s *Server) DefaultRespHeaders(next http.Handler) (processed http.Handler) {
	return
}

// WebsiteContentHandler - Handles all website content displayed when unknown/undesired clients
// connect to or ping the HTTP C2 server.
func (s *Server) WebsiteContentHandler(resp http.ResponseWriter, req *http.Request) {
}

func default404Handler(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(404)
}

// HandleRSAKey - Verifies the authenticity of the RSA key used by a ghost implant to prove its identity + encryption.
func (s *Server) HandleRSAKey(resp http.ResponseWriter, req *http.Request) {

}
