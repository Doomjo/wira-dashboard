package main

import (
    "database/sql"
    "log"
    "net/http"
    _ "github.com/lib/pq"
)

var db *sql.DB

func main() {
    // Database connection
    connStr := "postgres://postgres:123@localhost:5432/postgres?sslmode=disable"
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Test the connection
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

	insertTestData()
    // Initialize router with database connection
    router := SetupRouter(db)

    // Start server
    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}
