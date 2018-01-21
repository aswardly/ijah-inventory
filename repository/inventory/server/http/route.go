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
	addSKURoute.Handler(addSKUHandler)

	//updateSKU Route
	updateSKURoute := s.router.Path("/updateSKU")
	updateSKURoute.Methods("POST")
	serviceObj, found = s.sc.GetService("updateSKUHandler")
	if false == found {
		panic("service 'updateSKUHandler' not found")
	}
	updateSKUHandler, ok := serviceObj.(*handler.UpdateSKUHandler)
	if false == ok {
		panic("failed asserting 'updateSKUHandler'")
	}
	updateSKURoute.Handler(updateSKUHandler)

	//createSale Route
	createSaleRoute := s.router.Path("/createSale")
	createSaleRoute.Methods("POST")
	serviceObj, found = s.sc.GetService("createSaleHandler")
	if false == found {
		panic("service 'createSalehandler' not found")
	}
	createSaleHandler, ok := serviceObj.(*handler.CreateSaleHandler)
	if false == ok {
		panic("failed asserting 'createSaleHandler'")
	}
	createSaleRoute.Handler(createSaleHandler)

	//updateSale Route
	updateSaleRoute := s.router.Path("/updateSale")
	updateSaleRoute.Methods("POST")
	serviceObj, found = s.sc.GetService("updateSaleHandler")
	if false == found {
		panic("service 'updateSaleHandler' not found")
	}
	updateSaleHandler, ok := serviceObj.(*handler.UpdateSaleHandler)
	if false == ok {
		panic("failed asserting 'updateSaleHandler'")
	}
	updateSaleRoute.Handler(updateSaleHandler)

	//getAllStockValue Route
	getAllStockValueRoute := s.router.Path("/getStockValue")
	getAllStockValueRoute.Methods("GET")
	serviceObj, found = s.sc.GetService("getAllStockValueHandler")
	if false == found {
		panic("service 'getAllStockValueHandler' not found")
	}
	getAllStockValueHandler, ok := serviceObj.(*handler.GetAllStockValueHandler)
	if false == ok {
		panic("failed asserting 'getAllStockValueHandler'")
	}
	getAllStockValueRoute.Handler(getAllStockValueHandler)

	//getAllSalesValue Route
	getAllSalesValueRoute := s.router.Path("/getSalesValue")
	getAllSalesValueRoute.Methods("GET")
	serviceObj, found = s.sc.GetService("getAllSalesValueHandler")
	if false == found {
		panic("service 'getAllSalesValueHandler' not found")
	}
	getAllSalesValueHandler, ok := serviceObj.(*handler.GetAllSalesValueHandler)
	if false == ok {
		panic("failed asserting 'getAllStockValueHandler'")
	}
	getAllSalesValueRoute.Handler(getAllSalesValueHandler)

	//exportStockCSV route
	exportStockCSVRoute := s.router.Path("/exportStockCSV")
	exportStockCSVRoute.Methods("GET")
	serviceObj, found = s.sc.GetService("exportStockCSVHandler")
	if false == found {
		panic("service 'exportStockCSVHandler' not found")
	}
	exportStockCSVHandler, ok := serviceObj.(*handler.ExportStockCSVHandler)
	if false == ok {
		panic("failed asserting 'exportStockCSVHandler'")
	}
	exportStockCSVRoute.Handler(exportStockCSVHandler)
}
