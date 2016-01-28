package auth

import (
	"log"
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
