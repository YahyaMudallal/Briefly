package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ArticlesHandler handles HTTP requests for articles.
type ArticlesHandler struct {
	service *services.ArticleService
}

// NewArticlesHandler creates a new ArticlesHandler with dependency injection.
func NewArticlesHandler(service *services.ArticleService) *ArticlesHandler {
	return &ArticlesHandler{service: service}
}

// HandleGetArticles returns all articles.
func (h *ArticlesHandler) HandleGetArticles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := h.service.GetAllArticles(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve articles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// HandleGetArticle returns a single article by ID.
func (h *ArticlesHandler) HandleGetArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	article, err := h.service.GetArticleByID(ctx, id)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// HandleCreateArticle creates a new article.
func (h *ArticlesHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// create the article using the service layer
	ctx := r.Context()
	createdArticle, err := h.service.CreateArticle(ctx, &article)
	if err != nil {
		http.Error(w, "Failed to create article", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdArticle)
}

// HandleDeleteArticle deletes an article by ID.
func (h *ArticlesHandler) HandleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// call service layer to delete the article
	ctx := r.Context()
	err = h.service.DeleteArticle(ctx, id)
	if err != nil {
		http.Error(w, "Failed to delete article", http.StatusInternalServerError)
		return
	}

	// return no content status
	w.WriteHeader(http.StatusNoContent)
}