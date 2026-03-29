package services

import (
	"context"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
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
