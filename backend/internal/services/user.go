package services

import (
	"context"
	"fmt"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// UserService provides business logic for users.
type UserService struct {
	repository repositories.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repository repositories.UserRepository) *UserService {
	return &UserService{repository: repository}
}

// GetAllUsers retrieves all users.
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repository.GetAll(ctx)
}

// GetUserByID retrieves a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id bson.ObjectID) (*models.User, error) {
	return s.repository.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by their email address.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repository.GetByEmail(ctx, email)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, newUser *models.User) (*models.User, error) {

	// check if the email is valid
	if newUser.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// check if a user with the same email already exists
	existingUser, err := s.repository.GetByEmail(ctx, newUser.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", newUser.Email)
	}

	// check if the user first name is valid
	if newUser.FirstName == "" {
		return nil, fmt.Errorf("first name cannot be empty")
	}

	// check if the user last name is valid
	if newUser.LastName == "" {
		return nil, fmt.Errorf("last name cannot be empty")
	}

	// check if the password is valid
	if len(newUser.Password) < 4 {
		return nil, fmt.Errorf("password must be at least 4 characters long")
	}
	
	// hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	newUser.Password = string(hashedPassword)
	
	// set the admin status to false
	newUser.IsAdmin = false

	// set the created and updated timestamps
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	return s.repository.Create(ctx, newUser)
}