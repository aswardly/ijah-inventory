package http

import (
	"ijah-inventory/repository/inventory/server/http/handler"
)

//routeSetup is a function for setting up http routes
func (s *Server) routeSetup() {
	//register http routes here
	//Note: it is recommended to panic in the case of any http handler initialization failure (so there's no possibility of running a http server with broken http handler)
	//e.g. when getting a specific http handler as a named service object from service container and the service container reports no such service exists

	indexRoute := s.router.Path("/")
	indexRoute.Methods("GET")
	indexRoute.Handler(handler.NewDummyHandler(s.sc))
}
