package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/api"
	"github.com/YahyaMudallal/newsWebSite/internal/database"
	"github.com/joho/godotenv"
)

func main() {

	// lookup for a .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file has been found.")
	} else {
		log.Println(".env has been found.")
	}

	// get port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// get MongoDB URI and database name from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// get database name from environment variables
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "main"
	}

	// create a context with timeout for database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// initialize database connection
	db, err := database.NewDatabase(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer db.Close(context.Background())	// close the database connection when the application exits

	// define the application
	app := &api.Application{
		Database: db,
		Config: api.Config{
			Address: ":" + port,
		},
	}

	// Start the server
	err = app.Run(app.Mount())
	if err != nil {
		log.Fatal("Error at the start of the server:", err)
	}
}
