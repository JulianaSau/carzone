package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var db *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("Connecting to the database...")
	time.Sleep(5 * time.Second)

	// open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	// check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalf("error closing database: %v", err)
	}
}
