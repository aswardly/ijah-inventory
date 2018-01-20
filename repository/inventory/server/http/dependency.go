package http

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"ijah-inventory/repository/inventory/domain/inventory/datamapper"
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
	salesDatamapper := datamapper.NewPurchase(dbSession)
	s.sc.RegisterService("salesDatamapper", salesDatamapper)

	//test handler
	testHandler := &handler.TestHandler{}
	testHandler.SetContainer(s.sc)
	testHandler.Handle = testHandler.TestHandle
	s.sc.RegisterService("testHandler", testHandler)

	//perform injection
	if err := s.sc.Ready(); err != nil {
		panic(fmt.Sprintf("Service initialization failed with error: %+v", err))
	}
}
