package main

import (
	"flag"
	"fmt"
	"git"
	"net/http"
	"router"
	"strconv"
)

var connConf git.Config
var authFlag string
var portFlag int

func init() {
	flag.StringVar(&authFlag, "auth", "", "Git Basic Authentification key")
	flag.IntVar(&portFlag, "port", 2506, "Webserver port")
}

func main() {
	flag.Parse()
	
	connConf = git.Config{
		GitUrl:           "https://api.github.com",
		BaseOrganisation: "/informationgrid",
		GitAuthkey:       ""}

	if authFlag != "" {
		connConf.GitAuthkey = authFlag
	}

	git.SetConfig(connConf)
	
	RunCHM()

}

type Config struct {
	Port int
}

func RunCHM() {

	var config = Config{}
	
	config.Port = portFlag

	fmt.Println("starting CHM...")

	router := router.NewRouter()

	fmt.Println("starting server on Port", strconv.Itoa(config.Port))
	panic(http.ListenAndServe(":"+strconv.Itoa(config.Port), router))
}
