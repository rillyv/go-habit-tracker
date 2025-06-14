package main

import (
	"log"

	"github.com/rillyv/habit-tracker/db"
)

func main() {
	err := db.Migrate()

	if err != nil {
		log.Fatal("")
	}
}
