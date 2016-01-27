package controllers

import (
	"encoding/json"
	"expvar"
	"net/http"
)

/*  All util functions for API Calls
    Typically used to send JSON replies back to the client
    Fixed format calls for 200,400 & 500 status codes
*/

// RouteHits : Map for number of route hits
var RouteHits = expvar.NewMap("routeHits").Init()

// BasicResponse : JSON reply for API Calls
type BasicResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// RespondOrThrowErr : Respond to general requests or exit with server err.
func RespondOrThrowErr(responseObj BasicResponse, w http.ResponseWriter) {

	responseJSON, err := json.Marshal(responseObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// ThrowInternalErrAndExit : Respond with internal server error
func ThrowInternalErrAndExit(w http.ResponseWriter) {

	responsecontent := BasicResponse{
		"Internal Server Error",
		500,
	}

	w.WriteHeader(http.StatusInternalServerError)
	RespondOrThrowErr(responsecontent, w)
}

// ThrowForbiddenedAndExit : Used for requests whose resource is not found
func ThrowForbiddenedAndExit(w http.ResponseWriter) {

	responsecontent := BasicResponse{
		"Forbidden",
		403,
	}

	w.WriteHeader(http.StatusForbidden)
	RespondOrThrowErr(responsecontent, w)
}

// RespondSuccessAndExit : Repond with a success
func RespondSuccessAndExit(w http.ResponseWriter, msg string) {

	responsecontent := BasicResponse{
		msg,
		200,
	}
	w.WriteHeader(http.StatusOK)
	RespondOrThrowErr(responsecontent, w)

}
