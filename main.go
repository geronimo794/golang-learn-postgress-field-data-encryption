package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Data to be inserted
	username := "johndoe"
	email := "johndoe@example.com"
	passphrase := os.Getenv("ENCRYPTION_PASSPHRASE")
	if passphrase == "" {
		log.Fatal("ENCRYPTION_PASSPHRASE environment variable is not set")
	}

	// SQL statement with pgp_sym_encrypt
	query := `
        INSERT INTO users (username, email)
        VALUES ($1, pgp_sym_encrypt($2, $3))
    `

	// Execute the query
	_, err = db.Exec(query, username, email, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data inserted successfully.")

	getData(db, passphrase, username)
}

func getData(db *sql.DB, passphrase string, username string) {
	// SQL statement to retrieve and decrypt data
	query := `
    SELECT id, username, pgp_sym_decrypt(email, $1) AS email
    FROM users
    WHERE username = $2
`

	// Execute the query
	row := db.QueryRow(query, passphrase, username)

	// Variables to hold the retrieved data
	var id int
	var decryptedEmail string

	// Scan the result into variables
	err := row.Scan(&id, &username, &decryptedEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No user found with the given username.")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("User ID: %d\nUsername: %s\nEmail: %s\n", id, username, decryptedEmail)
	}
}
