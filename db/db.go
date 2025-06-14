package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() (*sql.DB, error) {
	fmt.Println("Trying to connect to the database...")

	connStr := "postgres://postgres:postgres@localhost/postgres?sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	fmt.Println("Connected to the database.")

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return nil, err
	}

	err = DB.Ping()

	if err != nil {
		fmt.Println("Failed to connect to the database.")
		return nil, err
	}

	return DB, nil
}

func Disconnect() error {
	err := DB.Close()

	if err != nil {
		fmt.Println("Failed to close the database.")
		return err
	}

	fmt.Println("Successfully closed the database.")
	return nil
}

func Migrate() error {
	driver, err := postgres.WithInstance(DB, &postgres.Config{})

	if err != nil {
		fmt.Println("Unable to create driver.")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)

	if err != nil {
		fmt.Println("Failed to find the migrations file")
		return err
	}

	err = m.Up()

	if err != nil {
		fmt.Println("Unable to migrate up.")
		return err
	}

	return nil
}
