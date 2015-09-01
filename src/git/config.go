package git

import (
	"strings"
	"time"
)

var config Config

func init() {
	config = Config{
		GitUrl:            "https://api.github.com",
		BaseOrganisation:  "informationgrid",
		GitAuthkey:        "",
		MiscDefaultBranch: "",
		MaxRepos:          100,
		MaxBranches:       100,
	}

	twoMonthAgo := time.Now().AddDate(0, -2, 0)
	config.SinceTime = twoMonthAgo.Format("2006-01-02")
}

type Config struct {
	GitUrl, BaseOrganisation, GitAuthkey, SinceTime, MiscDefaultBranch string
	MaxRepos, MaxBranches                                              int
}

func GetConfig() Config {
	return config
}


func SetConfig(connData Config) (settingsChanged bool) {
	if connData.GitUrl != "" && connData.GitUrl != config.GitUrl {
		config.GitUrl = connData.GitUrl
		settingsChanged = true
	}
	if connData.BaseOrganisation != "" && config.BaseOrganisation != connData.BaseOrganisation {
		config.BaseOrganisation = connData.BaseOrganisation
		settingsChanged = true
	}
	if strings.Replace(connData.GitAuthkey, "*", "", -1) != "" && config.GitAuthkey != connData.GitAuthkey {
		config.GitAuthkey = connData.GitAuthkey
		settingsChanged = true
	}
	if connData.SinceTime != "" && config.SinceTime != connData.SinceTime {
		config.SinceTime = connData.SinceTime
		settingsChanged = true
	}
	if connData.MaxRepos != 0 && config.MaxRepos != connData.MaxRepos {
		config.MaxRepos = connData.MaxRepos
		settingsChanged = true
	}
	if connData.MaxBranches != 0 && config.MaxBranches != connData.MaxBranches {
		config.MaxBranches = connData.MaxBranches
		settingsChanged = true
	}
	if connData.MiscDefaultBranch != "" && config.MiscDefaultBranch != connData.MiscDefaultBranch {
		config.MiscDefaultBranch = connData.MiscDefaultBranch
		settingsChanged = true
	}
	return settingsChanged
}

