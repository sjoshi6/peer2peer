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

// Chat is used to store chat details
type Chat struct {
	MentorID string `json:"mentorid"`
	MenteeID string `json:"menteeid"`
	ChatText string `json:"chattext"`
}

// Chats is an array of visitor
type Chats []ChatResponse

// ChatResponse is used to send a response to the client
type ChatResponse struct {
	Chat
	ID           string `json:"id"`
	CreationTime string `json:"creationtime"`
}

// ChatsResponse is used to send an arr of chats
type ChatsResponse struct {
	Chats Chats `json:"chats"`
}

// GetAllChatsRequest is used to take the mentor id input
type GetAllChatsRequest struct {
	MentorID string `json:"mentorid"`
}

// Create : Used to create a new visitor from http post request
func (c Chat) Create(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Create:POST /v1/chat", 1)

	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&c)
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
	stmt, _ := dbconn.Prepare(`INSERT INTO Chat(mentorid, menteeid, chattext) VALUES($1,$2,$3);`)

	_, execerr := stmt.Exec(c.MentorID, c.MenteeID, c.ChatText)

	if execerr != nil {
		// If execution err occurs then throw error
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	// If no error then give a success response
	RespondSuccessAndExit(w, "Chat Added Successfully")
}

// Get : Used to get data about the visitor
func (c Chat) Get(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Get:GET /v1/chat/{id}", 1)

	// Vars to extract values from DB ; necessary due to uneven struct
	var (
		chatid       string
		mentorid     string
		menteeid     string
		chattext     string
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
		QueryRow("SELECT id, mentorid, menteeid, chattext, creationtime FROM Chat WHERE id = $1", id).
		Scan(&chatid, &mentorid, &menteeid, &chattext, &creationtime)

	if err != nil {
		log.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	chatResponse := ChatResponse{
		Chat{mentorid, menteeid, chattext},
		chatid,
		creationtime.Format(time.RFC3339)}

	jsonResponse, err := json.Marshal(chatResponse)
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
func (c Chat) GetAll(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Get:GET /v1/visitors", 1)

	// Vars to extract values from DB ; necessary due to uneven struct
	var (
		mid          string
		chatid       string
		mentorid     string
		menteeid     string
		chattext     string
		creationtime time.Time
	)

	var getAllReq GetAllChatsRequest
	decoder := json.NewDecoder(r.Body)

	// Expand the json attached in post request
	err := decoder.Decode(&getAllReq)
	if err != nil {

		log.Println(err)
		ThrowForbiddenedAndExit(w)
		return
	}

	mid = getAllReq.MentorID

	// Used for per user connection to DB
	dbconn, err := db.GetDBConn(config.DBName)
	defer dbconn.Close()

	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	rows, err := dbconn.Query("SELECT * FROM Chat WHERE mentorid=$1", mid)
	if err != nil {
		log.Println(err)
		ThrowInternalErrAndExit(w)
		return
	}

	var chatArr Chats

	for rows.Next() {

		// Scan all the values for the given row
		rows.Scan(&chatid, &mentorid, &menteeid, &chattext, &creationtime)

		// Create a visitor object
		chatResponse := ChatResponse{
			Chat{mentorid, menteeid, chattext},
			chatid,
			creationtime.Format(time.RFC3339)}

		chatArr = append(chatArr, chatResponse)
	}

	chatsResp := ChatsResponse{
		Chats: chatArr,
	}

	jsonResponse, err := json.Marshal(chatsResp)
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
func (c Chat) Delete(w http.ResponseWriter, r *http.Request) {

	// Incr the debug vals
	RouteHits.Add("Delete:DELETE /v1/chat/{id}", 1)

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

	stmt, _ := dbconn.Prepare(`DELETE FROM Chat WHERE id=$1`)
	_, execerr := stmt.Exec(id)

	if execerr != nil {

		log.Println(execerr)
		ThrowInternalErrAndExit(w)
		return
	}

	RespondSuccessAndExit(w, "Chat Deleted SuccessFully")

}
