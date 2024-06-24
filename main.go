package main

import (
	"gitmate/config"
	"gitmate/routes"
	"log"
)

func main() {
	// Initialize the database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Set up routes
	router := routes.SetupRouter(db)
	router.Run(":8080")
}
