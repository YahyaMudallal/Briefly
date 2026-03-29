package handlers

import (
	"encoding/json"
	"net/http"

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
