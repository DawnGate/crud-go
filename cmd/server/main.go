package main

import (
	"log"
	"net/http"

	"sos-crud/internal/database"
	"sos-crud/internal/handlers"
	"sos-crud/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()

	r := mux.NewRouter()

	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUserByID).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.DeleteUserByID).Methods("DELETE")
	r.HandleFunc("/users/{id}", handlers.UpdateUserByID).Methods("PUT")

	r.HandleFunc("/sos", handlers.GetAllSos).Methods("GET")
	r.HandleFunc("/sos/{id}", handlers.GetSosByID).Methods("GET")

	r.HandleFunc("/signup", handlers.Signup).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)

	protected.HandleFunc("/profile", handlers.Profile).Methods("GET")

	protected.HandleFunc("/sos", handlers.CreateSos).Methods("POST")
	protected.HandleFunc("/sos/{id}", handlers.UpdateSosByID).Methods("PUT")
	protected.HandleFunc("/sos/{id}", handlers.DeleteSosByID).Methods("DELETE")

	log.Println("Server is running on port: 8080...")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Fail to start server: %v", err)
	}
}
