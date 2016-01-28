package controllers

import (
	"net/http"
	"peer2peer/auth"
	"peer2peer/credentials"

	"github.com/dgrijalva/jwt-go"
)

// Auth is a blank stuct used to namespace auth routes
type Auth struct{}

// Start a new Auth Server
var secretKey = credentials.SecretKey
var jwtProvider = auth.NewAuthManager([]byte(secretKey),
	auth.Config{Method: jwt.SigningMethodHS256, TTL: 3600 * 24 * 3})

// GrantToken : Route Handler to grant new token
func (a Auth) GrantToken(w http.ResponseWriter, req *http.Request) {

	// create a new JWT with claims, jwts adds "iat" and "exp" claims
	token := jwtProvider.NewToken()

	// This ID is to be set to the users ID recieved from the attached json not header.
	token.Claims["id"] = "sjoshi"
	tokenBytes, err := jwtProvider.SignToken(token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenBytes)
}
