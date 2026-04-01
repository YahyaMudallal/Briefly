package services

import (
	"context"

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
