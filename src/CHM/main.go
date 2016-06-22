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
var serverConfigFlag string
var port string
var host string

func init() {
	flag.StringVar(&serverConfigFlag, "conf", "server.cfg", "server.cfg path")
	flag.StringVar(&port, "port", "", "Webserver port (if set overwrites server.cfg config)")
	flag.StringVar(&host, "host", "", "Webserver host IP (if set overwrites server.cfg config)")

}

func main() {
	flag.Parse()

	err := router.LoadServerConfig(serverConfigFlag)
	if err != nil {
		log.Panicln(err)
	}
	if port == "" {
		port = router.GlobalServerConfig.Port
	}
	if host == "" {
		host = router.GlobalServerConfig.HostIP
	}
	router.InitJWT()

	RunCHM()

}

// RunCHM starts the programm
func RunCHM() {

	fmt.Println("starting CHM...")

	router := router.NewRouter()

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("starting server on:", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
