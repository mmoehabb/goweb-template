//go:build ignore

package main

import (
	"log"

	"goweb/db"
	"goweb/db/users"
)

func main() {
	if _, err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.RunMigrations(users.DataModel{}); err != nil {
		log.Fatalf("Failed to run GORM migrations: %v", err)
	}

	if err := db.RunGooseMigrations(); err != nil {
		log.Fatalf("Failed to run goose migrations: %v", err)
	}

	if err := db.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully.")
}
