package https

import "net/http"

// Handler - Path mapped to a handler function
type Handler func(resp http.ResponseWriter, req *http.Request)

// StartSession - An implant is trying to connect and register to server.
// This function handles all HTTP-specific details, before handing over the registration
// process to the 'c2' package, with HandleGhostRegistration() function.
func (s *Server) StartSession(resp http.ResponseWriter, req *http.Request) {

}

// HandleSession - The implant is registered: handle all possible requests/responses going from/to it.
func (s *Server) HandleSession(resp http.ResponseWriter, req *http.Request) {

}

// PollSession - Used when a user wants to have direct access to a remote system shell, with appropriate speed/latency.
func (s *Server) PollSession(resp http.ResponseWriter, req *http.Request) {

}

// StopSession - Kill the ghost implant connection.
func (s *Server) StopSession(resp http.ResponseWriter, req *http.Request) {

}

// GetSession - Returns an HTTP implant session.
func (s *Server) GetSession(req *http.Request) (sess *Session) {
	return
}
