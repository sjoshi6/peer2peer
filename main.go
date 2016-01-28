package main

import (
	"log"
	"net/http"
	"os"
	"peer2peer/db/postgres"
	"peer2peer/routers"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

}

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("Usage : ./peer2peer <PortNumber>")
	}

	port := ":" + string(args[0])

	log.Println("Go API Server - Logs", time.Now())
	log.Printf("API Server started at - %s\n", port)

	// Create all the missing tables before starting server
	db.CreateTablesIfNotExists()

	// Start the API Server
	StartServer(port)
}

// StartServer : Start the API Server by calling this function
func StartServer(port string) {

	// Creating a new mux router
	var router = mux.NewRouter().StrictSlash(true)

	// Add APP routes for various controllers
	router = routers.AddVisitorRoutes(router)

	// Add Auth routes
	router = routers.AddAuthRoutes(router)

	// This route is essential to view the monitoring stats for the app.
	router.Handle("/debug/vars", http.DefaultServeMux)

	log.Fatal(http.ListenAndServe(port, router))
}
