package main

import (
	"log"
	"net/http"
	"sos-crud/internal/database"
	"sos-crud/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")

	log.Println("Server is running on port: 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Fail to start server: %v", err)
	}
}
