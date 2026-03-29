package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/services"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// UsersHandler handles HTTP requests for users.
type UsersHandler struct {
	service *services.UserService
}

// NewUsersHandler creates a new UsersHandler with dependency injection.
func NewUsersHandler(service *services.UserService) *UsersHandler {
	return &UsersHandler{service: service}
}

// HandleGetUsers returns all users.
func (h *UsersHandler) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.service.GetAllUsers(ctx)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// HandleGetUser returns a single user by ID.
func (h *UsersHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
