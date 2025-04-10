package main

import (
	"log"

	"github.com/hoodnoah/ghoam/cmd/migrate"
)

func main() {
	_, err := migrate.Execute("data/ghoam.db")
	if err != nil {
		log.Fatalf("failed to initialize the database with error %v", err)
	}

	log.Println("successfully initialized the database")
}
