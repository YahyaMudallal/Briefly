package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/api"
	"github.com/YahyaMudallal/newsWebSite/internal/clients"
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

	newAPIKey := os.Getenv("NEWS_API_KEY")
	if newAPIKey == "" {
		newAPIKey = "PLACEHOLDER_API_KEY"
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		geminiAPIKey = "PLACEHOLDER_API_KEY"
	}

	// create a context with timeout for database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// initialize database connection
	db, err := database.NewDatabase(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer db.Close(context.Background()) // close the database connection when the application exits

	// Initialize the external news API client
	newClient := clients.NewNewsAPIClient(newAPIKey)

	// Initialize the external Gemini API client
	geminiClient := clients.NewGeminiAPIClient(geminiAPIKey) 

	// define the application
	app := &api.Application{
		Database: db,
		NewsClient: newClient,
		GeminiClient: geminiClient,
		Config: api.Config{
			Address: ":" + port,
		},
	}

	// Mount the handlers and get the main handler
	handler := app.Mount()

	// Start the scheduler for daily article synchronization
	app.Scheduler.Start()
	defer app.Scheduler.Stop() // stop the scheduler when the application exits

	// Run the server
	err = app.Run(handler)
	if err != nil {
		log.Fatal("Error at the start of the server:", err)
	}
}
