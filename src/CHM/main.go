package main

import (
	"fmt"
	"git"
	"net/http"
	"router"
	"strconv"
)

var connConf git.Config

func init() {
	connConf = git.Config{
		GitUrl:           "https://api.github.com",
		BaseOrganisation: "/informationgrid",
		GitAuthkey:       "QlJvZWhybDpZcCFtSzZGMw=="}

	git.SetConfig(connConf)
}

func main() {

	RunCHM()

}

type Config struct {
	Port int
}

func RunCHM() {

	var config = Config{
		Port: 2506,
	}

	fmt.Println("starting CHM...")

	router := router.NewRouter()

	fmt.Println("starting server on Port", strconv.Itoa(config.Port))
	panic(http.ListenAndServe(":"+strconv.Itoa(config.Port), router))
}
