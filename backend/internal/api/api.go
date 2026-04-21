package api

import (
	"log"
	"net/http"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/clients"
	"github.com/YahyaMudallal/newsWebSite/internal/database"
	"github.com/YahyaMudallal/newsWebSite/internal/handlers"
	"github.com/YahyaMudallal/newsWebSite/internal/jobs"
	"github.com/YahyaMudallal/newsWebSite/internal/middleware"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"github.com/YahyaMudallal/newsWebSite/internal/services"
)

// Define collection names as constants
const articlesCollection = "articles"
const usersCollection = "users"
const commentsCollection = "comments"

// Application represents the main structure of the server.
// Contains configuration, database, and dependencies.
type Application struct {
	Config   Config
	Database *database.Database
	NewsClient clients.NewsClient
	Scheduler *jobs.Scheduler
	GeminiClient clients.GeminiClient
}

// Config represents the configuration of the server.
type Config struct {
	Address string
}

// Mount defines the routes of the server and returns a ServeMux with the handlers.
func (app *Application) Mount() http.Handler {
	mux := http.NewServeMux()

	// Create repositories
	articlesRepo := repositories.NewMongoArticleRepository(
		app.Database.GetCollection(articlesCollection),
	)
	usersRepo := repositories.NewMongoUserRepository(
		app.Database.GetCollection(usersCollection),
	)
	commentsRepo := repositories.NewMongoCommentRepository(
		app.Database.GetCollection(commentsCollection),
	)

	// Create services
	articlesService := services.NewArticleService(articlesRepo, usersRepo, commentsRepo, app.NewsClient, app.GeminiClient)
	usersService := services.NewUserService(usersRepo)
	commentsService := services.NewCommentService(commentsRepo, usersRepo)

	// Create handlers with dependency injection
	articlesHandler := handlers.NewArticlesHandler(articlesService)
	commentsHandler := handlers.NewCommentsHandler(commentsService)
	usersHandler := handlers.NewUsersHandler(usersService)

	// Initialize the scheduler for daily article synchronization
	app.Scheduler = jobs.NewScheduler(articlesService)

	// Define the routes and their handlers

	mux.HandleFunc("GET /v1/articles", articlesHandler.HandleGetArticles)
	mux.HandleFunc("GET /v1/articles/{id}", articlesHandler.HandleGetArticle)
	mux.HandleFunc("GET /v1/comments", commentsHandler.HandleGetComments)
	mux.HandleFunc("GET /v1/comments/{id}", commentsHandler.HandleGetComment)
	mux.HandleFunc("GET /v1/users", usersHandler.HandleGetUsers)
	mux.HandleFunc("GET /v1/users/{id}", usersHandler.HandleGetUser)
	mux.HandleFunc("GET /v1/comments/article/{articleId}", commentsHandler.HandleGetCommentsByArticle)
	mux.HandleFunc("POST /v1/users", usersHandler.HandleCreateUser)
	mux.HandleFunc("POST /v1/users/login", usersHandler.HandleLoginUser)
	mux.HandleFunc("POST /v1/articles", middleware.AuthMiddleware(articlesHandler.HandleCreateArticle))
	mux.HandleFunc("POST /v1/comments", middleware.AuthMiddleware(commentsHandler.HandleCreateComment))
	mux.HandleFunc("DELETE /v1/comments/{id}", middleware.AuthMiddleware(commentsHandler.HandleDeleteComment))
	mux.HandleFunc("DELETE /v1/articles/{id}", middleware.AuthMiddleware(articlesHandler.HandleDeleteArticle))
	mux.HandleFunc("DELETE /v1/users/{id}", middleware.AuthMiddleware(usersHandler.HandleDeleteUser))

	//first wrapping the router with CORS
	corsHandler := middleware.EnableCORS(mux)
	// Adding middlewares
	logHandler := middleware.LoggingMiddleware(corsHandler)
	recoverHandler := middleware.RecoverMiddleware(logHandler)

	return recoverHandler
}

// Run starts the http server.
// Returns an error if the server hasn't started.
func (app *Application) Run(handler http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Address,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server started at %s\n", app.Config.Address)

	return srv.ListenAndServe()
}
