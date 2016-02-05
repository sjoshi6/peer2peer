package main

import (
	"log"
	"net/http"
	"os"
	"peer2peer/config"
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
	if len(args) < 1 {
		log.Fatal("Usage : ./peer2peer <mode>")
	}

	mode := args[0]

	// Create all the missing tables before starting server
	db.CreateTablesIfNotExists()

	switch mode {

	case "auth":
		log.Println("Auth Server - Logs", time.Now())
		log.Printf("Auth Server started at - %s\n", config.AuthServerPort)

		// Start the Auth Server
		StartAuthServer(config.AuthServerPort)

	case "api":
		log.Println("API Server - Logs", time.Now())
		log.Printf("API Server started at - %s\n", config.APIServerPort)

		// Start the API Server
		StartServer(config.APIServerPort)

	default:
		log.Fatal("Incorrect Choice of mode : Only auth / api is legal")

	}

}

// StartServer : Start the API Server by calling this function
func StartServer(port string) {

	// Need to activate this for token based access
	var a controllers.Auth

	// Creating a new mux router
	var router = mux.NewRouter().StrictSlash(true)

	// Add APP routes for various controllers
	router = routers.AddVisitorRoutes(router)

	// This route is essential to view the monitoring stats for the app.
	router.Handle("/debug/vars", http.DefaultServeMux)

	n := negroni.Classic()

	// Need to activate this for token based access
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
