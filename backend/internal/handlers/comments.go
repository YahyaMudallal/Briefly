package handlers

import (
	"encoding/json"
	"net/http"

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
	ctx := r.Context()
	comments, err := h.service.GetAllComments(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// HandleGetComment returns a single comment by ID.
func (h *CommentsHandler) HandleGetComment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	comment, err := h.service.GetCommentByID(ctx, id)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
