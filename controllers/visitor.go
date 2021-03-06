package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"peer2peer/config"
	"peer2peer/db/postgres"
	"time"

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

// Visitors is an array of visitor
type Visitors []VisitorResponse

// VisitorResponse : Used to reply back to a visitor get request
type VisitorResponse struct {
	Visitor
	ID           string `json:"id" db:"id"`
	CreationTime string `json:"creationtime" db:"creationtime"`
}

// VisitorsResponse is used to send back a reply of visitors array
type VisitorsResponse struct {
	Visitors Visitors `json:"visitors"`
}

// Create : Used to create a new visitor from http post request
func (v Visitor) Create(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Create:POST /v1/visitor", 1)
	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&v)
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

	// Insert into DB
	stmt, _ := dbconn.Prepare(`INSERT INTO Visitor(firstname, lastname, age, gender, email,
                               phonenumber, university) VALUES($1,$2,$3,$4,$5,$6,$7);`)

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

	// Incr the debug vals
	RouteHits.Add("Get:GET /v1/visitor/{id}", 1)

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
		creationtime time.Time
	)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ThrowForbiddenedAndExit(w)
		return
	}

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	err = dbconn.
		QueryRow("SELECT id, firstname, lastname, age, gender, email, phonenumber,university, creationtime FROM Visitor WHERE id = $1", id).
		Scan(&visitorid, &fname, &lname, &age, &gender, &email, &phonenumber, &university, &creationtime)

	if err != nil {
		log.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	visitorResp := VisitorResponse{
		Visitor{fname, lname, age, gender, email, phonenumber, university},
		visitorid,
		creationtime.Format(time.RFC3339)}

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

// GetAll is used to get all the visitors
func (v Visitor) GetAll(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Get:GET /v1/visitors", 1)

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
		creationtime time.Time
	)

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	rows, err := dbconn.Query("SELECT * FROM Visitor")
	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	var visitorArr Visitors

	for rows.Next() {

		// Scan all the values for the given row
		rows.Scan(&visitorid, &fname, &lname, &age, &gender, &email, &phonenumber, &university, &creationtime)

		// Create a visitor object
		visitorResp := VisitorResponse{
			Visitor{fname, lname, age, gender, email, phonenumber, university},
			visitorid,
			creationtime.Format(time.RFC3339)}

		visitorArr = append(visitorArr, visitorResp)
	}

	visitorsResp := VisitorsResponse{
		Visitors: visitorArr,
	}

	jsonResponse, err := json.Marshal(visitorsResp)
	if err != nil {

		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	// Append the data to response writer
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

// Delete : Delete a visitor from DB
func (v Visitor) Delete(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Delete:DELETE /v1/visitor/{id}", 1)

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		ThrowForbiddenedAndExit(w)
		return
	}

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {

		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	stmt, _ := dbconn.Prepare(`DELETE FROM Visitor WHERE id=$1`)
	_, execerr := stmt.Exec(id)

	if execerr != nil {

		log.Println(execerr)
		ThrowInternalErrAndExit(w)
		return
	}

	RespondSuccessAndExit(w, "Visitor Deleted SuccessFully")

}
