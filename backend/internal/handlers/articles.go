package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
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
	
	// default values for pagination
	page := 1
	limit := 3
	
	// read URL query parameters (?page=1&limit=10)
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}
	
	// read URL query parameters for sorting
	sortBy := r.URL.Query().Get("sortBy")
    if sortBy == "" {
        sortBy = "date"
    }
    order := r.URL.Query().Get("order")
    if order == "" {
        order = "desc"
    }

	articles, err := h.service.GetPaginated(ctx, page, limit, sortBy, order)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
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
		http.Error(w, err.Error(), apperrors.FilterError(err))
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
		http.Error(w, err.Error(), apperrors.FilterError(err))
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
	articleID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// get the id of the user from the context
	userID, ok := r.Context().Value("user_id").(bson.ObjectID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// call service layer to delete the article
	ctx := r.Context()
	err = h.service.DeleteArticle(ctx, articleID, userID)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}

	// return no content status
	w.WriteHeader(http.StatusNoContent)
}

// HandleGenerateSummary generates a summary for an article by ID.
func (h *ArticlesHandler) HandleGenerateSummary(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	articleID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}
  
 	// call the service layer to generate the summary
	ctx := r.Context()
	summary, err := h.service.GenerateSummary(ctx, articleID)
  if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}
  
  	// return the generated summary
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

// HandleUpvoteArticle upvotes an article by ID.
func (h *ArticlesHandler) HandleToggleUpvote(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	articleID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}
  
	// get the id of the user from the context
	userID, ok := r.Context().Value("user_id").(bson.ObjectID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// call service layer to upvote the article
	ctx := r.Context()
	nbUpvotes, err := h.service.ToggleUpvote(ctx, articleID, userID)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}
  
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"upvotes": nbUpvotes})
}

// HandleDownvoteArticle downvotes an article by ID.
func (h *ArticlesHandler) HandleToggleDownvote(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	articleID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// get the id of the user from the context
	userID, ok := r.Context().Value("user_id").(bson.ObjectID)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// call service layer to downvote the article
	ctx := r.Context()
	nbDownvotes, err := h.service.ToggleDownvote(ctx, articleID, userID)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"downvotes": nbDownvotes})
}
