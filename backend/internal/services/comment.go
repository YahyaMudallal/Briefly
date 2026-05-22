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
	articleRepo repositories.ArticleRepository
	commentRepo repositories.CommentRepository
	userRepo    repositories.UserRepository
}

// NewCommentService creates a new CommentService.
func NewCommentService(
	commentRepo repositories.CommentRepository,
	userRepo repositories.UserRepository,
	articleRepo repositories.ArticleRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		userRepo: userRepo,
		articleRepo: articleRepo}
}

// GetAllComments retrieves all comments.
func (s *CommentService) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	return s.commentRepo.GetAll(ctx)
}

// GetCommentByID retrieves a comment by its ID.
func (s *CommentService) GetCommentByID(ctx context.Context, id bson.ObjectID) (*models.Comment, error) {
	return s.commentRepo.GetByID(ctx, id)
}

// GetCommentsByArticleID retrieves all comments for a specific article.
func (s *CommentService) GetCommentsByArticleID(ctx context.Context, articleID bson.ObjectID) ([]models.Comment, error) {
	return s.commentRepo.GetByArticleID(ctx, articleID)
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {

	// check if the content is empty
	if comment.Content == "" {
		return nil, fmt.Errorf("%w : content cannot be empty", apperrors.ErrValidation)
	}

	// check if the article exists
	_, err := s.articleRepo.GetByID(ctx, comment.ArticleID)
	if err != nil {
		return nil, fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	// check if the user exists
	_, err = s.userRepo.GetByID(ctx, comment.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
	}

	// increment the number of comments for the article
	err = s.articleRepo.IncrementCommentCount(ctx, comment.ArticleID, 1)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to update article's comment count: %w", apperrors.ErrInternal, err)
	}

	// set the created at and updated at fields
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	return s.commentRepo.Create(ctx, comment)
}

// DeleteComment deletes a comment by its ID.
func (s *CommentService) DeleteComment(ctx context.Context, commentID bson.ObjectID, userID bson.ObjectID) error {

	// check if the comment exists
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("%w : comment not found", apperrors.ErrNotFound)
	}

	// check if the user is the author of the comment
	isAuthor := comment.AuthorID == userID

	// check if the user is an admin
	isAdmin := false
	if !isAuthor {
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
		}
		isAdmin = user.IsAdmin
	}

	// check if the user have the permission to delete the comment
	if !isAuthor && !isAdmin {
		return fmt.Errorf("%w : unauthorized to delete this comment", apperrors.ErrUnauthorized)
	}

	// change the comment to be deleted
	comment.Content = "[deleted]"
	comment.AuthorID = bson.NilObjectID
	comment.UpdatedAt = time.Now()

	// Update the comment in the database
	return s.commentRepo.Update(ctx, comment)
}
