package main

import (
	"log"
	"net/http"
	"os"
	"peer2peer/controllers"
	"peer2peer/db/postgres"
	"peer2peer/routers"
	"runtime"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

}

func main() {

	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Usage : ./peer2peer <mode> <PortNumber>")
	}

	mode := args[0]
	port := ":" + string(args[1])

	switch mode {

	case "auth":
		log.Println("Auth Server - Logs", time.Now())
		log.Printf("Auth Server started at - %s\n", port)
		StartAuthServer(port)

	case "api":
		log.Println("API Server - Logs", time.Now())
		log.Printf("API Server started at - %s\n", port)

		// Create all the missing tables before starting server
		db.CreateTablesIfNotExists()

		// Start the API Server
		StartServer(port)

	default:
		log.Fatal("Incorrect Choice of mode : Only auth / api is legal")

	}

}

// StartServer : Start the API Server by calling this function
func StartServer(port string) {

	var a controllers.Auth

	// Creating a new mux router
	var router = mux.NewRouter().StrictSlash(true)

	// Add APP routes for various controllers
	router = routers.AddVisitorRoutes(router)

	// This route is essential to view the monitoring stats for the app.
	router.Handle("/debug/vars", http.DefaultServeMux)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(a.RequireTokenAuthentication))
	n.UseHandler(router)
	n.Run(port)

}

// StartAuthServer is used to start the auth server
func StartAuthServer(port string) {

	// Creating a new mux router
	var authrouter = mux.NewRouter().StrictSlash(true)

	// Add Auth routes
	authrouter = routers.AddAuthRoutes(authrouter)

	n := negroni.Classic()

	n.UseHandler(authrouter)
	n.Run(port)
}
