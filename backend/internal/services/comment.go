package services

import (
	"context"
	"fmt"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// CommentService provides business logic for comments.
type CommentService struct {
	repository repositories.CommentRepository
}

// NewCommentService creates a new CommentService.
func NewCommentService(repository repositories.CommentRepository) *CommentService {
	return &CommentService{repository: repository}
}

// GetAllComments retrieves all comments.
func (s *CommentService) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	return s.repository.GetAll(ctx)
}

// GetCommentByID retrieves a comment by its ID.
func (s *CommentService) GetCommentByID(ctx context.Context, id bson.ObjectID) (*models.Comment, error) {
	return s.repository.GetByID(ctx, id)
}

// GetCommentsByArticleID retrieves all comments for a specific article.
func (s *CommentService) GetCommentsByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Comment, error) {
	return s.repository.GetByArticleID(ctx, articleID)
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	// check if the content is empty
	if comment.Content == "" {
		return nil, fmt.Errorf("%w : content cannot be empty", apperrors.ErrValidation)
	}

	// check if the article exists
	_, err := s.repository.GetByArticleID(ctx, comment.ArticleID)
	if err != nil {
		return nil, fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	// set the created at and updated at fields
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return s.repository.Create(ctx, comment)
}

// DeleteComment deletes a comment by its ID.
func (s *CommentService) DeleteComment(ctx context.Context, id bson.ObjectID) error {

	// check if the comment exists
	_, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("%w : comment not found", apperrors.ErrNotFound)
	}

	return s.repository.Delete(ctx, id)
}
