package git

import (
	"strings"
	"time"
)

var config Config

func init() {
	config = Config{
		GitURL:            "https://api.github.com",
		BaseOrganisation:  "informationgrid",
		GitAuthkey:        "",
		MiscDefaultBranch: "",
		MaxRepos:          100,
		MaxBranches:       100,
	}

	config.SinceTime = time.Now().AddDate(0, -2, 0)
}

// Config contains all settings which can be saved
type Config struct {
	GitURL, BaseOrganisation, GitAuthkey, MiscDefaultBranch string
	MaxRepos, MaxBranches                                   int
	SinceTime                                               time.Time
}

// GetConfig returns the current config
func GetConfig() Config {
	return config
}

// SetConfig updates the current Config to connData.
// The return values indicate if the default Branch changed
// and if all Data must be reloaded to use the updated Settings
func SetConfig(connData Config) (updateAll, miscBranchChanged bool) {
	if connData.GitURL != "" && connData.GitURL != config.GitURL {
		config.GitURL = connData.GitURL
		updateAll = true
	}
	if connData.BaseOrganisation != "" && config.BaseOrganisation != connData.BaseOrganisation {
		config.BaseOrganisation = connData.BaseOrganisation
		updateAll = true
	}
	if strings.Replace(connData.GitAuthkey, "*", "", -1) != "" && config.GitAuthkey != connData.GitAuthkey {
		config.GitAuthkey = connData.GitAuthkey
		updateAll = true
	}
	if !(connData.SinceTime.Equal(config.SinceTime) || connData.SinceTime.Equal(time.Time{})) {
		config.SinceTime = connData.SinceTime
		updateAll = false
	}
	if connData.MaxRepos != 0 && config.MaxRepos != connData.MaxRepos {
		config.MaxRepos = connData.MaxRepos
		updateAll = true
	}
	if connData.MaxBranches != 0 && config.MaxBranches != connData.MaxBranches {
		config.MaxBranches = connData.MaxBranches
		updateAll = true
	}
	if connData.MiscDefaultBranch != "" && config.MiscDefaultBranch != connData.MiscDefaultBranch {
		config.MiscDefaultBranch = connData.MiscDefaultBranch
		updateAll = false
		miscBranchChanged = true
	}
	return
}
