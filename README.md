# commit-finder
Commit finder is a webapp for searching through all repositories of an organisation simultaneously.


## Build from source

Prerequisites:

* [Git](https://git-scm.com/downloads) - `brew install git`, `apt-get install git`, etc
* [Go 1.6 or higher](https://golang.org/dl/).

Clone Repository

    git clone https://github.com/BRoehrl/commit-finder

Set GOPATH and get dependencies

    export GOPATH=$YOUR_DIRECTORY/commit-finder
    go get github.com/gorilla/mux
    go get github.com/dgrijalva/jwt-go

Install commit-finder

    go install CHM

Register a new application at https://github.com/settings/developers with the intended server adress as the callback URL and write down the Client ID and Client Secret.


Create a `server.cfg` file matching this template:

    {
    	"GitClientID": "THE_CLIENT_ID",
    	"GitClientSecret": "THE_MATCHING_CLIENT_SECRET",
    	"JWTSecret": "SOME_SECRET_FOR_JWT_SIGNING",
    	"HostIP": "127.0.0.1",
    	"Port": "8080"
    }

Run the webapp

`./bin/CHM` or `./bin/CHM.exe`
