package git

var userIDTokens map[string](string)

func init() {
	userIDTokens = make(map[string](string))
}

// AddUser adds or updates a user in userIDTokens
func AddUser( userID, secret string){
  userIDTokens[userID] = secret
}

// GetAccessToken returns the AccessToken and true or "" and false if userID
// wasn't found
func GetAccessToken(userID string) (string, bool){
  secret, ok := userIDTokens[userID]
  return secret, ok
}
