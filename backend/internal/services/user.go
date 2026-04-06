package services

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
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

	// check if the email isn't empty
	if newUser.Email == "" {
		return nil, fmt.Errorf("%w : email cannot be empty", apperrors.ErrValidation)
	}

	// check if the email is in a valid format
	_, err := mail.ParseAddress(newUser.Email)
	if err != nil {
		return nil, fmt.Errorf("%w : invalid email format", apperrors.ErrValidation)
	}

	// check if a user with the same email already exists
	existingUser, err := s.repository.GetByEmail(ctx, newUser.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("%w : user with email %s already exists", apperrors.ErrConflict, newUser.Email)
	}

	// check if the user first name is valid
	if newUser.FirstName == "" {
		return nil, fmt.Errorf("%w : first name cannot be empty", apperrors.ErrValidation)
	}

	// check if the user last name is valid
	if newUser.LastName == "" {
		return nil, fmt.Errorf("%w : last name cannot be empty", apperrors.ErrValidation)
	}

	// check if the password is valid
	if len(newUser.Password) < 4 {
		return nil, fmt.Errorf("%w : password must be at least 4 characters long", apperrors.ErrValidation)
	}
	
	// hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w : failed to hash password", apperrors.ErrInternal)
	}
	newUser.Password = string(hashedPassword)
	
	// set the admin status to false
	newUser.IsAdmin = false

	// set the created and updated timestamps
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	return s.repository.Create(ctx, newUser)
}

// LoginUser authenticates a user and return the user model if it is sucessful, or an error if it fails.
func (s *UserService) LoginUser(ctx context.Context, email string, password string) (*models.User, error) {

	// check if the email isn't empty
	if email == "" {
		return nil, fmt.Errorf("%w : email cannot be empty", apperrors.ErrValidation)
	}

	// check if the email is in a valid format
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("%w : invalid email format", apperrors.ErrValidation)
	}

	// check if the password is valid
	if password == "" {
		return nil, fmt.Errorf("%w : password cannot be empty", apperrors.ErrValidation)
	}

	// get the user by email
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("%w : invalid email or password", apperrors.ErrUnauthorized)
	}

	// compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w : invalid email or password", apperrors.ErrUnauthorized)
	}

	return user, nil
}