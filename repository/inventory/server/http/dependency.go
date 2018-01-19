package http

import (
	"fmt"
	"strconv"

	httpConfig "ijah-inventory/repository/inventory/server/config/http"
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

	//perform injection
	if err := s.sc.Ready(); err != nil {
		panic(fmt.Sprintf("Service initialization failed with error: %+v", err))
	}
}
