package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// HandleCreateUser creates a new user and generates a JWT token.
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
		http.Error(w, err.Error(), apperrors.FilterError(err))
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

// HandleLoginUser authenticates a user and returns a JWT token.
func (h *UsersHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	
	// create a temporary struct to hold the incoming JSON data
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// authenticate the user using the service layer
	ctx := r.Context()
	loggedUser, err := h.service.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}

	// remove the password
	loggedUser.Password = ""

	// generate a JWT token for the user
	tokenString, err := auth.GenerateToken(loggedUser.ID)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to generate JWT token: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// prepare the response
	response := map[string]interface{}{
		"user": loggedUser,
		"token": tokenString,
	}

	// send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleDeleteUser delete the user.
func (h *UsersHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {

	// parse query
	idStr := r.PathValue("id")
	deletedUserID, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// get the id of the user from the context
	userID, ok := r.Context().Value("user_id").(bson.ObjectID)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// call service layer
	ctx := r.Context()
	err = h.service.DeleteUser(ctx, deletedUserID, userID)
	if err != nil {
		http.Error(w, err.Error(), apperrors.FilterError(err))
		return
	}

	// return success response
	w.WriteHeader(http.StatusNoContent)
}