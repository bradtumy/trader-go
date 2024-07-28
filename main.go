package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	// Load configuration
	var cfg Config
	readFile(&cfg)

	// Initialize database connection
	err := initializeDatabase(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Router setup
	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/users", returnAllUsers).Methods("GET")
	router.HandleFunc("/createUser", createUser).Methods("POST")
	router.HandleFunc("/stocks", returnAllStocks).Methods("GET")
	router.HandleFunc("/orders", returnAllOrders).Methods("GET")
	router.HandleFunc("/createOrder", createOrder).Methods("POST")
	router.HandleFunc("/createStock", createStock).Methods("POST")
	router.HandleFunc("/user/{id}/profile", GetUserProfile).Methods("GET")
	router.HandleFunc("/searchStocks", searchStocks).Methods("GET")
	router.HandleFunc("/keyMetrics", getKeyMetrics).Methods("GET")

	fmt.Println("Starting server on port", cfg.Server.Port)
	log.Fatal(http.ListenAndServeTLS(cfg.Server.Host+":"+cfg.Server.Port, cfg.Server.Certificate, cfg.Server.Key, router))
}
