package services

import (
	"context"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ArticleService provides business logic for articles.
type ArticleService struct {
	repository repositories.ArticleRepository
}

// NewArticleService creates a new ArticleService.
func NewArticleService(repository repositories.ArticleRepository) *ArticleService {
	return &ArticleService{repository: repository}
}

// GetAllArticles retrieves all articles.
func (s *ArticleService) GetAllArticles(ctx context.Context) ([]models.Article, error) {
	return s.repository.GetAll(ctx)
}

// GetArticleByID retrieves an article by its ID.
func (s *ArticleService) GetArticleByID(ctx context.Context, id bson.ObjectID) (*models.Article, error) {
	return s.repository.GetByID(ctx, id)
}
