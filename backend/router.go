// router.go
package main

import (
    "database/sql"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func SetupRouter(db *sql.DB) http.Handler {
    router := mux.NewRouter()

    // Create new handler instance
    handler := NewHandler(db)

    // Define routes
    router.HandleFunc("/api/players", handler.GetPlayers).Methods("GET")

    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:8080"),
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    })

    return c.Handler(router)
}