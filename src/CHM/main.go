package main

import (
	"flag"
	"fmt"
	"git"
	"log"
	"net/http"
	"router"
)

var connConf git.Config
var authFlag string
var portFlag int
var host string


func init() {
	flag.StringVar(&authFlag, "auth", "", "Git Basic Authentification key")
	flag.IntVar(&portFlag, "port", 2506, "Webserver port")
	flag.StringVar(&host, "host", "127.0.0.1", "Webserver host IP")
	
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

	addr := fmt.Sprintf("%s:%d", host, config.Port)
	fmt.Println("starting server on:", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
