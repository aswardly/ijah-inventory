//Package database is for database related configurations
package database

//Config is a collection of configuration items
type Config struct {
	DbFile string //path to database file
}

//StartUp allows the config to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (c *Config) StartUp() {
	//initialize the startup process here
}

//Shutdown allows the config to satisfy gocontainer.Service interface (import package github.com/ncrypthic/gocontainer)
func (c *Config) Shutdown() {
	//perform any shutdown process here
}
