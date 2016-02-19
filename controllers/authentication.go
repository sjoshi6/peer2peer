package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"peer2peer/auth"
	"peer2peer/config"
	"peer2peer/credentials"
	"peer2peer/db/postgres"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

var hashCost = config.Cost

// Auth is a blank stuct used to namespace auth routes
type Auth struct{}

// SignUp is used to handle signup requests
type SignUp struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	AdminToken string `json:"admintoken"`
}

// Login is used to handle login requests
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Start a new Auth Server
var secretKey = credentials.SecretKey
var jwtProvider = auth.NewAuthManager([]byte(secretKey),
	auth.Config{Method: jwt.SigningMethodHS256, TTL: 3600 * 24 * 3})

// grantToken is a local function used to grant new tokens to the user
func (a Auth) grantToken(userid string) ([]byte, error) {

	// create a new JWT with claims, jwts adds "iat" and "exp" claims
	token := jwtProvider.NewToken()

	// This ID is to be set to the users ID recieved from the attached json not header.
	token.Claims["id"] = "sjoshi6"
	tokenBytes, err := jwtProvider.SignToken(token)

	log.Printf("Granting a new token to user %s \n", string(tokenBytes))

	if err != nil {
		log.Println("Could not generate token")
		return nil, err
	}

	return tokenBytes, nil

}

// RequireTokenAuthentication is used to ensure request token is validated before serving
func (a Auth) RequireTokenAuthentication(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	if req.Method == "OPTIONS" {
        log.Println("Found options skip check")
		next(w, req)
	}

	token, err := jwtProvider.GetToken(req)

	if err == nil && token.Valid {

		log.Println("Valid token found")
		log.Println(token.Claims["id"])
		next(w, req)

	} else {

		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse, err := json.Marshal(BasicResponse{"Unauthorized", 401})

		if err != nil {
			log.Println(err)
		}

		w.Write(jsonResponse)
	}
}

// SignUpHandler is used to manage signup requests
func (a Auth) SignUpHandler(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("POST /v1/signup", 1)

	// Creating a reference to hold signup data
	var su SignUp

	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&su)
	if err != nil {
		log.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	if su.AdminToken != config.AdminToken {

		log.Println("UI AdminToken did not match")
		ThrowForbiddenedAndExit(w)
		return

	}

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {
		ThrowInternalErrAndExit(w)
		return
	}

	// Insert statement prepare
	stmt, _ := dbconn.Prepare(`INSERT INTO SignUp(email, password) VALUES($1,$2);`)

	hash, err := bcrypt.GenerateFromPassword([]byte(su.Password), hashCost)
	if err != nil {

		log.Println("bcrypt hash creation broke")
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return

	}

	// if no hash err then execute statement

	_, execerr := stmt.Exec(su.Email, string(hash))
	if execerr != nil {

		// If execution err occurs then throw error
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	// If no error then give a success response
	// Aquire a new token to send the user
	token, err := a.grantToken("sjoshi6")

	if err != nil {

		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(token)

}

// LoginHandler is used to handle login requests
func (a Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var dbPassword string

	// Incr the debug vals
	RouteHits.Add("POST /v1/login", 1)

	// Creating a reference to hold signup data
	var loginReq Login

	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&loginReq)
	if err != nil {
		log.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {
		ThrowInternalErrAndExit(w)
		return
	}

	err = dbconn.
		QueryRow("SELECT password from SignUp WHERE email= $1", loginReq.Email).Scan(&dbPassword)

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	loginerr := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(loginReq.Password))
	if loginerr != nil {

		// If err is thrown credentials are mismatched
		responsecontent := BasicResponse{
			"Login Credentials are incorrect",
			400,
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Status", "Client Error")

		RespondOrThrowErr(responsecontent, w)
		return
	}

	// If no error in comparehash means login Credentials match

	// Aquire a new token to send the user
	token, err := a.grantToken("sjoshi6")

	if err != nil {

		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(token)

}

//OptionsHandler is used to handle the options route
func (a Auth) OptionsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write([]byte("Done"))
}
