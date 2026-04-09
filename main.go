package main

import (
	"embed"
	"log"

	"github.com/guiflauzino18/economizze/internal/infra/database"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {

	if err := database.RunMigrations(nil, migrationsFS); err != nil {
		log.Fatal("migrations")
	}

}
