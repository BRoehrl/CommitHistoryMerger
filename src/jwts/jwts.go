// Copied from https://github.com/dghubble/jwts
// Author: Dalton Hubble

// Package jwts offers a Provider interface for JSON Web Token management backed
// by github.com/dgrijalva/jwt-go JWTs.
package jwts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const defaultTTL = 3600 * 24 * 7 // 1 week

// A Provider creates, signs, and retrieves JSON Web Tokens (JWTs).
type Provider interface {
	// New returns a new JWT Token using the Store's signing method.
	New() *jwt.Token
	// Sign digitally signs the Token to return the JWT byte slice.
	Sign(token *jwt.Token) ([]byte, error)
	// Get gets the valid JWT from the Authorization header. If the token is
	// missing, expired, or the signature does not validate, returns an error.
	Get(req *http.Request) (*jwt.Token, error)
}

// Config configures a Manager.
type Config struct {
	// digital signing method, defaults to jwt.SigningMethodHS256 (SHA256)
	Method jwt.SigningMethod
	// token expiration time in seconds, defaults to 1 week
	TTL int64
}

// Manager is a JSON Web Token (JWT) Provider which create or retrieves tokens
// with a particular signing key and options.
type Manager struct {
	key    []byte
	method jwt.SigningMethod
	ttl    int64
}

// New creates a new Manager which provides JWTs using the given signing key.
// Defaults to signing with SHA256 HMAC (jwt.SigningMethodHS256) and expiring
// tokens after 1 week.
func New(key []byte, configs ...Config) *Manager {
	var c Config
	if len(configs) == 0 {
		c = Config{}
	} else {
		c = configs[0]
	}
	m := &Manager{
		key:    key,
		method: c.Method,
		ttl:    c.TTL,
	}
	m.setDefaults()
	return m
}

func (m *Manager) setDefaults() {
	if m.method == nil {
		m.method = jwt.SigningMethodHS256
	}
	if m.ttl == 0 {
		m.ttl = defaultTTL
	}
}

// New returns a new *jwt.Token which has the prescribed signing method, issued
// at time, and expiration time set on it.
//
// Add claims to the Claims map and use the controller to Sign(token) to get
// the standard JWT signed string representation.
func (m *Manager) New() *jwt.Token {
	token := jwt.New(m.method)
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	claims["iatStr"] = fmt.Sprintf("%v", time.Now().Unix())
	d := time.Duration(m.ttl) * time.Second
	claims["exp"] = time.Now().Add(d).Unix()
	return token
}

// Sign digitally signs a *jwt.Token using the token's method and the manager's
// signing key to return a JWT byte slice.
func (m *Manager) Sign(token *jwt.Token) ([]byte, error) {
	jwtString, err := token.SignedString(m.key)
	return []byte(jwtString), err
}

// Get gets the signed JWT from the Authorization header or jwt cookie. If the token is
// missing, expired, or the signature does not validate, returns an error.
func (m *Manager) Get(req *http.Request) (*jwt.Token, error) {

	jwtString := req.Header.Get("Authorization")
	if jwtString == "" {
		cookie, err := req.Cookie("jwt")
		if err != nil {
			return nil, request.ErrNoTokenInRequest
		}
		if cookie.Value == "" {
			return nil, request.ErrNoTokenInRequest
		}
		jwtString = cookie.Value
	}
	token, err := jwt.Parse(jwtString, m.getKey)
	if err == nil && token.Valid {
		// token parsed, exp/nbf checks out, signature verified, Valid is true
		return token, nil
	}
	return nil, request.ErrNoTokenInRequest
}

// getKey accepts an unverified JWT and returns the signing/verification key.
// Also ensures tha the token's algorithm matches the signing method expected
// by the manager.
func (m *Manager) getKey(unverified *jwt.Token) (interface{}, error) {
	// require token alg to match the set signing method, do not allow none
	if meth := unverified.Method; meth == nil || meth.Alg() != m.method.Alg() {
		return nil, jwt.ErrHashUnavailable
	}
	return m.key, nil
}
