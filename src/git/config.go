package git

import (
	"strings"
	"time"
)

// DefaultConfig is the default configuration
var DefaultConfig Config

func init() {
	DefaultConfig = GetNewConfig()
}

// Config contains all settings which can be saved
type Config struct {
	GitURL,
	BaseOrganisation,
	GitAuthkey,
	MiscDefaultBranch,
	JWTissuedAt string
	MaxRepos,
	MaxBranches int
	SinceTime time.Time
}

// GetNewConfig returns a new Config with default values
func GetNewConfig() Config {
	c := Config{}
	c.GitURL = "https://api.github.com"
	c.BaseOrganisation = "informationgrid"
	c.MiscDefaultBranch = "develop"
	c.MaxRepos = 100
	c.MaxBranches = 100
	c.SinceTime = time.Now().AddDate(0, -6, 0)
	return c
}

// GetDefaultConfig returns the default config
func GetDefaultConfig() Config {
	return DefaultConfig
}

// SetConfig updates the userConfig to connData.
// The return values indicate if the default Branch changed
// and if all Data must be reloaded to use the updated Settings
func SetConfig(userConfig *Config, connData Config) (updateAll, miscBranchChanged bool) {
	if connData.GitURL != "" && connData.GitURL != userConfig.GitURL {
		userConfig.GitURL = connData.GitURL
		updateAll = true
	}
	if connData.BaseOrganisation != "" && userConfig.BaseOrganisation != connData.BaseOrganisation {
		userConfig.BaseOrganisation = connData.BaseOrganisation
		updateAll = true
	}
	if strings.Replace(connData.GitAuthkey, "*", "", -1) != "" && userConfig.GitAuthkey != connData.GitAuthkey {
		userConfig.GitAuthkey = connData.GitAuthkey
		updateAll = true
	}
	if !(connData.SinceTime.Equal(userConfig.SinceTime) || connData.SinceTime.Equal(time.Time{})) {
		userConfig.SinceTime = connData.SinceTime
		updateAll = false
	}
	if connData.MaxRepos != 0 && userConfig.MaxRepos != connData.MaxRepos {
		userConfig.MaxRepos = connData.MaxRepos
		updateAll = true
	}
	if connData.MaxBranches != 0 && userConfig.MaxBranches != connData.MaxBranches {
		userConfig.MaxBranches = connData.MaxBranches
		updateAll = true
	}
	if connData.MiscDefaultBranch != "" && userConfig.MiscDefaultBranch != connData.MiscDefaultBranch {
		userConfig.MiscDefaultBranch = connData.MiscDefaultBranch
		updateAll = true
		miscBranchChanged = true
	}
	if connData.JWTissuedAt != "" {
		userConfig.JWTissuedAt = connData.JWTissuedAt
	}
	return
}
