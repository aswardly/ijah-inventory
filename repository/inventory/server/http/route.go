package http

import (
	"ijah-inventory/repository/inventory/server/http/handler"
)

//routeSetup is a function for setting up http routes
func (s *Server) routeSetup() {
	//register http routes here
	//Note: it is recommended to panic in the case of any http handler initialization failure (so there's no possibility of running a http server with broken http handler)
	//e.g. when getting a specific http handler as a named service object from service container and the service container reports no such service exists

	//index route
	indexRoute := s.router.Path("/")
	indexRoute.Methods("GET")
	indexRoute.Handler(handler.NewDummyHandler(s.sc))

	//test route
	testRoute := s.router.Path("/test")
	testRoute.Methods("GET")
	serviceObj, found := s.sc.GetService("testHandler")
	if false == found {
		panic("service 'testhandler' not found")
	}
	testHandler, ok := serviceObj.(*handler.TestHandler)
	if false == ok {
		panic("failed asserting 'testHandler'")
	}
	testRoute.Handler(testHandler)

	//getItemInfo Route
	getItemInfoRoute := s.router.Path("/itemInfo")
	getItemInfoRoute.Methods("GET")
	serviceObj, found = s.sc.GetService("getItemInfoHandler")
	if false == found {
		panic("service 'getItemInfoHandler' not found")
	}
	getItemInfoHandler, ok := serviceObj.(*handler.GetItemInfoHandler)
	if false == ok {
		panic("failed asserting 'getItemInfoHandler'")
	}
	getItemInfoRoute.Handler(getItemInfoHandler)

	//addSKU Route
	addSKURoute := s.router.Path("/addSKU")
	addSKURoute.Methods("POST")
	serviceObj, found = s.sc.GetService("addSKUHandler")
	if false == found {
		panic("service 'addSKUHandler' not found")
	}
	addSKUHandler, ok := serviceObj.(*handler.AddSKUHandler)
	if false == ok {
		panic("failed asserting 'addSKUHandler'")
	}
	getItemInfoRoute.Handler(addSKUHandler)
}
