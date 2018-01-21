package http

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"ijah-inventory/repository/inventory/domain/inventory/datamapper"
	"ijah-inventory/repository/inventory/domain/inventory/service"
	dbConfig "ijah-inventory/repository/inventory/server/config/database"
	httpConfig "ijah-inventory/repository/inventory/server/config/http"
	"ijah-inventory/repository/inventory/server/http/handler"
)

//setup is a function where setup of the http server is performed
func (s *Server) setup() {
	//set server options here (such as address, etc)
	serverPort := s.config.GetString("http.port")
	s.Addr = ":" + serverPort

	//register services to the service container here
	//e.g.
	//fooService := "foo service"
	//s.sc.RegisterService("foo", fooService)

	//compose http config object and register it as a service object
	httpPort, err := strconv.Atoi(s.config.GetString("http.port"))
	if err != nil {
		panic(fmt.Sprintf("Error converting string to int for http port config: %+v", err))
	}
	httpConfigObj := &httpConfig.Config{
		ListenAddress: s.config.GetString("http.address"),
		ListenPort:    httpPort,
		AccessLogPath: s.config.GetString("http.combinedLog.path"),
		AppLogPath:    s.config.GetString("http.appLog.path"),
	}
	s.sc.RegisterService("httpConfig", httpConfigObj)

	//database config
	databaseConfig := &dbConfig.Config{
		DbFile: s.config.GetString("database.filePath"),
	}
	//open db session
	dbSession, err := sql.Open("sqlite3", databaseConfig.DbFile)
	if err != nil {
		panic(fmt.Sprintf("Database initialization failed: %v", err))
	}

	//stock datamapper
	stockDatamapper := datamapper.NewStock(dbSession)
	s.sc.RegisterService("stockDatamapper", stockDatamapper)

	//purchase datamapper
	purchaseDatamapper := datamapper.NewPurchase(dbSession)
	s.sc.RegisterService("purchaseDatamapper", purchaseDatamapper)

	//sales datamapper
	salesDatamapper := datamapper.NewSale(dbSession)
	s.sc.RegisterService("salesDatamapper", salesDatamapper)

	//inventory service
	inventoryService := &service.Inventory{}
	s.sc.RegisterService("inventoryService", inventoryService)

	//test handler
	testHandler := &handler.TestHandler{}
	testHandler.SetContainer(s.sc)
	testHandler.Handle = testHandler.TestHandle
	s.sc.RegisterService("testHandler", testHandler)

	//getItemInfo Handler
	getItemInfoHandler := &handler.GetItemInfoHandler{}
	getItemInfoHandler.SetContainer(s.sc)
	getItemInfoHandler.Handle = getItemInfoHandler.GetItemInfoHandle
	s.sc.RegisterService("getItemInfoHandler", getItemInfoHandler)

	//addSKU Handler
	addSKUHandler := &handler.AddSKUHandler{}
	addSKUHandler.SetContainer(s.sc)
	addSKUHandler.Handle = addSKUHandler.AddSKUHandle
	s.sc.RegisterService("addSKUHandler", addSKUHandler)

	//updateSKU Handler
	updateSKUHandler := &handler.UpdateSKUHandler{}
	updateSKUHandler.SetContainer(s.sc)
	updateSKUHandler.Handle = updateSKUHandler.UpdateSKUHandle
	s.sc.RegisterService("updateSKUHandler", updateSKUHandler)

	//createSale Handler
	createSaleHandler := &handler.CreateSaleHandler{}
	createSaleHandler.SetContainer(s.sc)
	createSaleHandler.Handle = createSaleHandler.CreateSaleHandle
	s.sc.RegisterService("createSaleHandler", createSaleHandler)

	//updateSale Handler
	updateSaleHandler := &handler.UpdateSaleHandler{}
	updateSaleHandler.SetContainer(s.sc)
	updateSaleHandler.Handle = updateSaleHandler.UpdateSaleHandle
	s.sc.RegisterService("updateSaleHandler", updateSaleHandler)

	//getAllStockValue Handler
	getAllStockValueHandler := &handler.GetAllStockValueHandler{}
	getAllStockValueHandler.SetContainer(s.sc)
	getAllStockValueHandler.Handle = getAllStockValueHandler.GetAllStockValueHandle
	s.sc.RegisterService("getAllStockValueHandler", getAllStockValueHandler)

	//getAllSalesValue Handler
	getAllSalesValueHandler := &handler.GetAllSalesValueHandler{}
	getAllSalesValueHandler.SetContainer(s.sc)
	getAllSalesValueHandler.Handle = getAllSalesValueHandler.GetAllSalesValueHandle
	s.sc.RegisterService("getAllSalesValueHandler", getAllSalesValueHandler)

	//perform injection
	if err := s.sc.Ready(); err != nil {
		panic(fmt.Sprintf("Service initialization failed with error: %+v", err))
	}
}
