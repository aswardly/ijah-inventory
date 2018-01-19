//Package http is for http server related configurations
package http

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

//Config is a collection of configuration items
type Config struct {
	ListenAddress   string    //listening address of the http server
	ListenPort      int       //listening port no of the http server
	AccessLogPath   string    //path of the http access log
	AppLogPath      string    //path of the http application log path
	AccessLogWriter io.Writer //io writer for the http access log
	AppLogWriter    io.Writer //io writer for the http application log
}

//StartUp allows the config to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (c *Config) StartUp() {
	//initialize the io writers here
	accessLogFile, err := os.OpenFile(c.AccessLogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("HTTP config initialization, error opening access log file: %v", err))
	}
	c.AccessLogWriter = accessLogFile

	appLogFile, err := os.OpenFile(c.AppLogPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("HTTP config initialization, error opening app log file: %v", err))
	}
	c.AppLogWriter = appLogFile

	//setup global logger (formatter, output, level)
	log.SetFormatter(new(log.JSONFormatter))
	log.SetOutput(c.AppLogWriter)
	log.SetLevel(log.DebugLevel)
}

//Shutdown allows the config to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (c *Config) Shutdown() {
	//close the io writers (if they are closeable, i.e. implements io.Closer)
	accessLogCloser, ok := c.AccessLogWriter.(io.Closer)
	if ok {
		accessLogCloser.Close()
	}

	appLogCloser, ok := c.AppLogWriter.(io.Closer)
	if ok {
		appLogCloser.Close()
	}
}
