package main

import (
	"log"

	"github.com/hoodnoah/ghoam/cmd/migrate"
	"github.com/hoodnoah/ghoam/cmd/server"
)

func main() {
	// set up DB; run migrations and return a repository abstracting db interactions
	repos, err := migrate.Execute("data/ghoam.db")
	if err != nil {
		log.Fatalf("failed to initialize the database with error %v", err)
	}

	log.Println("successfully initialized the database")

	// set up webserver
	server.Execute(repos)
}
