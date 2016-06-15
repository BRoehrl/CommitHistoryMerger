package user

import (
	"git"
	"log"
)

var userIDTokens map[string](string)
var allUserCaches map[string](git.UserCache)
var userConfigLoaded map[string](bool)

func init() {
	userIDTokens = make(map[string](string))
	userConfigLoaded = make(map[string](bool))
	allUserCaches = make(map[string](git.UserCache))
}

// AddUser adds or updates a user in userIDTokens
func AddUser(userID, secret string) {
	userIDTokens[userID] = secret
	newCache := git.GetNewUserCache()
	newCache.UserID = userID
	newConfig := git.GetNewConfig()
	newConfig.GitAuthkey = secret
	newCache.Config = newConfig
	allUserCaches[userID] = newCache
}

// ConfigLoaded returns if the config of the user has been externally loaded
func ConfigLoaded(userID string) bool {
	return userConfigLoaded[userID]
}

// SetConfigLoaded sets if the config of the user has been externally loaded
func SetConfigLoaded(userID string, loaded bool) {
	userConfigLoaded[userID] = loaded
}

// GetAccessToken returns the AccessToken and true or "" and false if userID
// wasn't found
func GetAccessToken(userID string) (string, bool) {
	secret, ok := userIDTokens[userID]
	return secret, ok
}

// Exists checks if a User exists
func Exists(userID string) bool {
	_, ok := userIDTokens[userID]
	return ok
}

// GetUserCache returns the UserCache of userID
func GetUserCache(userID string) git.UserCache {
	userCache, ok := allUserCaches[userID]
	if !ok {
		log.Println("ERROR: UserCache from User '" + userID + "' not found")
	}
	return userCache
}

// SetUserCache sets the UserCache of userID
func SetUserCache(userID string, userCache git.UserCache) {
	allUserCaches[userID] = userCache
}

// SetCachedCommits sets the CachedCommits of userID
func SetCachedCommits(userID string, cachedCommits git.Commits) {
	userCache := allUserCaches[userID]
	userCache.CachedCommits = cachedCommits
	allUserCaches[userID] = userCache
}
