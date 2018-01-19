package main

import (
	"fmt"
	"path"
	"runtime"

	"ijah-inventory/repository/inventory/server/http"

	"github.com/gorilla/mux"
	"github.com/ncrypthic/gocontainer"
	"github.com/spf13/viper"
)

var sc *gocontainer.ServiceRegistry
var router *mux.Router

func main() {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if ok == false {
		panic("HTTP Server Run: can't get current file name")
	}

	//parse config file and compose the config objects
	httpConfigPath := path.Join(path.Dir(currentFilePath), "../../config/http")
	dbConfigPath := path.Join(path.Dir(currentFilePath), "../../config/database")

	config := viper.New()
	//http config
	config.SetConfigName("httpConfig")   //name of config file (without extension)
	config.AddConfigPath(httpConfigPath) //path to look for the config file in
	err := config.ReadInConfig()         //read the config file
	if err != nil {
		panic(fmt.Errorf("Failed reading http config: %v", err))
	}

	//db config
	config.SetConfigName("dbConfig")   //name of config file (without extension)
	config.AddConfigPath(dbConfigPath) //path to look for the config file in
	err = config.MergeInConfig()       //merge the config file
	if err != nil {
		panic(fmt.Errorf("Failed merging database config: %v", err))
	}

	//service container
	sc = gocontainer.NewContainer()

	//gorilla mux router for routing
	router = mux.NewRouter()

	//the http server
	server := http.NewServer(router, sc, config)
	server.Run()
}
