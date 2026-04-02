package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/auth"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
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

// HandleCreateUser creates a new user.
func (h* UsersHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {

	// create a temporary struct to hold the incoming JSON data
	var req struct {
		Email	string `json:"email"`
		FirstName string `json:"firstName"`
		LastName string `json:"lastName"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// transform the temporary struct into a User model
	user := models.User{
        Email:     req.Email,
        FirstName: req.FirstName,
        LastName:  req.LastName,
        Password:  req.Password,
    }

	// create the user using the service layer
	ctx := r.Context()
	createdUser, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		// filter the error based on its type and return the appropriate HTTP status code and message
		if errors.Is(err, services.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} 
		if errors.Is(err, services.ErrConflict) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// remove the password field from the response
	createdUser.Password = ""

	// generate a JWT token for the created user
	tokenString, err := auth.GenerateToken(createdUser.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to generate JWT token: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// prepare the response
	response := map[string]interface{}{
		"user": createdUser,
		"token": tokenString,
	}

	// send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}