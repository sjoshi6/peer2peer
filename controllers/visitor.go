package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"peer2peer/config"
	"peer2peer/db/postgres"

	"github.com/gorilla/mux"
)

// Visitor : A new student requesting mentorship
type Visitor struct {
	FirstName   string `json:"firstname" db:"lastname"`
	LastName    string `json:"lastname" db:"lastname"`
	Age         string `json:"age" db:"age"`
	Gender      string `json:"gender" db:"gender"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phonenumber" db:"phonenumber"`
	University  string `json:"university" db:"university"`
}

// VistorResponse : Used to reply back to a visitor get request
type VistorResponse struct {
	Visitor
	ID           string `json:"id" db:"id"`
	CreationTime string `json:"creationtime" db:"creationtime"`
}

// Create : Used to create a new visitor from http post request
func (v Visitor) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	// Insert into DB
	stmt, _ := dbconn.Prepare(`INSERT INTO Visitor(firstname, lastname, age, gender, email,
                               phonenumber, univeristy) VALUES($1,$2,$3,$4,$5,$6,$7);`)

	_, execerr := stmt.Exec(v.FirstName, v.LastName, v.Age,
		v.Gender, v.Email, v.PhoneNumber, v.University)

	if execerr != nil {
		// If execution err occurs then throw error
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	// If no error then give a success response
	RespondSuccessAndExit(w, "Visitor Added Successfully")
}

// Get : Used to get data about the visitor
func (v Visitor) Get(w http.ResponseWriter, r *http.Request) {

	// Vars to extract values from DB ; necessary due to uneven struct
	var (
		visitorid    string
		fname        string
		lname        string
		age          string
		gender       string
		email        string
		phonenumber  string
		university   string
		creationtime string
	)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ThrowForbiddenedAndExit(w)
		return
	}

	// Used for per user connection to DB
	dbconn := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	err := dbconn.
		QueryRow("SELECT id, firstname, lastname, age, gender, email, phonenumber,university, creationtime FROM users WHERE id = $1", id).
		Scan(&visitorid, &fname, &lname, &age, &gender, &email, &phonenumber, &university, &creationtime)

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	visitorResp := VistorResponse{
		Visitor{fname, lname, age, gender, email, phonenumber, university},
		visitorid,
		creationtime}

	jsonResponse, err := json.Marshal(visitorResp)
	if err != nil {
		ThrowInternalErrAndExit(w)
		return
	}

	// Append the data to response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
