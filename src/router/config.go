package router

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

// The ServerConfig struct contains all basic server options
type ServerConfig struct {
	GitClientID,
	GitClientSecret,
	JWTSecret,
	HostIP,
	Port string
}

// The GlobalServerConfig contains all information loaded on program startup
var GlobalServerConfig ServerConfig

// LoadServerConfig loads a Config-file into GlobalServerConfig
func LoadServerConfig(filePath string) (err error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Serverconfig not found: ", filePath)
		return
	}
	decoder := json.NewDecoder(file)
	serverConfig := ServerConfig{}
	err = decoder.Decode(&serverConfig)
	if err != nil {
		log.Println("Could not parse Config-file: ", file.Name())
		return
	}
	if serverConfig.hasEmptyFields() {
		err = errors.New("ServerConfig has empty fields. Please check the config file")
		return
	}
	GlobalServerConfig = serverConfig
	return nil
}

func (sc ServerConfig) hasEmptyFields() bool {
	if sc.GitClientID == "" {
		return true
	}
	if sc.GitClientSecret == "" {
		return true
	}
	if sc.JWTSecret == "" {
		return true
	}
	if sc.HostIP == "" {
		return true
	}
	if sc.Port == "" {
		return true
	}
	return false
}
