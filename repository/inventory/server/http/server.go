package http

import (
	"net/http"
	"os"

	httpConfig "ijah-inventory/repository/inventory/server/config/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ncrypthic/gocontainer"
	"github.com/spf13/viper"
)

//Server is a representation of http server
type Server struct {
	*http.Server                              //embedded http server object
	router       *mux.Router                  //gorilla mux router
	sc           gocontainer.ServiceContainer //service container
	config       *viper.Viper                 //http server config
}

//NewServer creates a new http server and returns a reference to it
func NewServer(rt *mux.Router, sc gocontainer.ServiceContainer, config *viper.Viper) *Server {
	return &Server{
		&http.Server{},
		rt,
		sc,
		config,
	}
}

//Run is a function for running the http server
func (s *Server) Run() {
	s.setup()
	s.routeSetup()

	//get http config object from service container (should have been registered during server setup)
	config, found := s.sc.GetService("httpConfig")
	if false == found {
		panic("Could not get http config from service container")
	}
	configObj, ok := config.(*httpConfig.Config)
	if false == ok {
		panic("Failed asserting http config as *httpConfig.Config")
	}

	_, ok = configObj.AccessLogWriter.(*os.File)
	if false == ok {
		panic("Failed asserting config.AccessLogWriter as *os.File")
	}

	//wrap the http server handler with gorilla/mux.CombinedLoggingHandler for apache style combined access log
	s.Handler = handlers.CombinedLoggingHandler(configObj.AccessLogWriter, s.router)
	s.ListenAndServe()
}
