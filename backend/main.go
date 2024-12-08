package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	// Add command line flag
	seedData := flag.Bool("seed", false, "Seed the database with test data")
	flag.Parse()

	if *seedData {
		// Use pgxpool for seeding
		connStr := "postgresql://postgres:123@localhost:5432/postgres"
		pool, err := pgxpool.Connect(context.Background(), connStr)
		if err != nil {
			log.Fatal(err)
		}
		defer pool.Close()

		log.Println("Starting data seeding...")
		err = insertTestData(context.Background(), pool)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Data seeding completed!")
		return
	}

	// Regular server startup
	connStr := "postgres://postgres:123@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")

	router := SetupRouter(db)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

