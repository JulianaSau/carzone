package driver

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
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

	// log snv vars
	fmt.Println("DB_HOST: ", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT: ", os.Getenv("DB_PORT"))
	fmt.Println("DB_USER: ", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD: ", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME: ", os.Getenv("DB_NAME"))

	// actual string
	// connStr := "host=db port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	fmt.Println("Connecting to the database...")
	var err error
	for i := 0; i < 10; i++ { // Retry up to 10 times
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("Successfully connected to the database")
			return
		}

		fmt.Printf("Database not ready (%d/10), retrying in 2 seconds: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Could not connect to the database after retries: %v", err)
	// open a connection to the database
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatalf("error opening database: %v", err)
	// }

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
	if err := db.Close(); err != nil {
		log.Fatalf("error closing database: %v", err)
	}
}
