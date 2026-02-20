package main

import (
	"fmt"
	"os"
)

func main() {
	dbURL := os.Getenv("QRAP_DATABASE_URL")
	if dbURL == "" {
		fmt.Fprintln(os.Stderr, "QRAP_DATABASE_URL is required")
		os.Exit(1)
	}

	fmt.Println("QRAP database migrations")
	fmt.Println("========================")
	fmt.Printf("Database: %s\n\n", dbURL)
	fmt.Println("Run migrations with:")
	fmt.Println("  migrate -path db/migrations -database $QRAP_DATABASE_URL up")
}
