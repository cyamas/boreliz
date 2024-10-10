package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cyamas/boreliz/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in the .env file")
	}
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	server.Run(conn)
}
