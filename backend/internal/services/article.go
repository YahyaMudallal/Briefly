package services

import (
	"context"
	"time"

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

// CreateArticle creates a new article.
func (s *ArticleService) CreateArticle(ctx context.Context, article *models.Article) (*models.Article, error) {

	// initialize the upvotes and downvotes to zero
	article.UpVotes = 0
	article.DownVotes = 0

	// initialize the summary to the empty string
	article.Summary = ""

	// set the created at and updated at fields
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	return s.repository.Create(ctx, article)
}

// DeleteArticle deletes an article by its ID.
func (s *ArticleService) DeleteArticle(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Delete(ctx, id)
}
