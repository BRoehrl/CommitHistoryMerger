package main

import (
	"flag"
	"fmt"
	"git"
	"git/processor"
	"log"
	"net/http"
	"router"
)

var connConf git.Config
var authFlag string
var configFlag string
var portFlag int
var host string

func init() {
	flag.StringVar(&authFlag, "auth", "", "Git Basic Authentification key")
	flag.StringVar(&configFlag, "conf", "default", "Use existing config (overwrittes auth flag)")
	flag.IntVar(&portFlag, "port", 2506, "Webserver port")
	flag.StringVar(&host, "host", "127.0.0.1", "Webserver host IP")

}

func main() {
	flag.Parse()

	err := processor.LoadCompleteConfig(configFlag)
	if err != nil {
		connConf = git.Config{
			GitUrl:            "https://api.github.com",
			BaseOrganisation:  "informationgrid",
			MiscDefaultBranch: "develop",
			GitAuthkey:        ""}

		if authFlag != "" {
			connConf.GitAuthkey = authFlag
		}

		git.SetConfig(connConf)
	}

	RunCHM()

}

// RunCHM starts the programm
func RunCHM() {

	fmt.Println("starting CHM...")

	router := router.NewRouter()

	addr := fmt.Sprintf("%s:%d", host, portFlag)
	fmt.Println("starting server on:", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
