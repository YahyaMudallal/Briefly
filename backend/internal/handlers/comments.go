package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CommentsHandler handles HTTP requests for comments.
type CommentsHandler struct {
	service *services.CommentService
}

// NewCommentsHandler creates a new CommentsHandler with dependency injection.
func NewCommentsHandler(service *services.CommentService) *CommentsHandler {
	return &CommentsHandler{service: service}
}

// HandleGetComments returns all comments.
func (h *CommentsHandler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	// call service layer
	ctx := r.Context()
	comments, err := h.service.GetAllComments(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return comments as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// HandleGetComment returns a single comment by ID.
func (h *CommentsHandler) HandleGetComment(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// call service layer
	ctx := r.Context()
	comment, err := h.service.GetCommentByID(ctx, id)
	if err != nil {
		// filter error
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return comment as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

// HandleGetCommentsByArticle returns all comments for a specific article.
func (h *CommentsHandler) HandleGetCommentsByArticle(w http.ResponseWriter, r *http.Request) {
	// parse query
	articleIDStr := r.PathValue("articleId")
	articleID, err := bson.ObjectIDFromHex(articleIDStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	// call service layer
	ctx := r.Context()
	comments, err := h.service.GetCommentsByArticleID(ctx, articleID)
	if err != nil {
		// filter error
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return comments as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}


// HandleCreateComment creates a new comment.
func (h *CommentsHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	// parse request body
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// call service layer to create the comment
	ctx := r.Context()
	createdComment, err := h.service.CreateComment(ctx, &comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the created comment as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdComment)
}

// HandleDeleteComment deletes a comment by ID.
func (h *CommentsHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	// parse query
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// call service layer
	ctx := r.Context()
	err = h.service.DeleteComment(ctx, id)
	if err != nil {
		// filter error
		if errors.Is(err, apperrors.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return success response
	w.WriteHeader(http.StatusNoContent)
}