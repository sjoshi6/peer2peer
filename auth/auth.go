package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const defaultTTL = 3600 * 24 * 7 // 1 week

// Config configures a JWT Manager.
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

// NewToken is a manager function. It takes no params and returns
// a new jwt.Token Object
func (m *Manager) NewToken() *jwt.Token {

	// expiry timer
	expTimer := time.Duration(m.ttl) * time.Second

	token := jwt.New(m.method)
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(expTimer).Unix()
	return token

}

// SignToken is used to sign the jwt token. It takes jwt.Token input and returns byte and err
func (m *Manager) SignToken(token *jwt.Token) ([]byte, error) {

	// Take the input token and sign it with manager key
	jwtStr, err := token.SignedString(m.key)

	if err != nil {
		log.Println(err)
	}

	return []byte(jwtStr), err
}

// GetToken gets the signed JWT from the Authorization header. If the token is
// missing, expired, or the signature does not validate, returns an error.
func (m *Manager) GetToken(req *http.Request) (*jwt.Token, error) {

	jwtString := req.Header.Get("Authorization")
	if jwtString == "" {

		// No auth header
		return nil, jwt.ErrNoTokenInRequest
	}

	// parse the jwt auth String with the manager key and get back a token
	// The getkey validates that the incoming token has valid methods and alogs.
	// If all is valid it returns back a key to the jwt.Parse function
	token, err := jwt.Parse(jwtString, m.getKey)

	if err == nil && token.Valid {
		// token parsed, exp/nbf checks out, signature verified, Valid is true
		return token, nil
	}
	return nil, jwt.ErrNoTokenInRequest
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
