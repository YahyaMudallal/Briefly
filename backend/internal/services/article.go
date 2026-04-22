package services

import (
	"context"
	"fmt"
	"time"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/clients"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"github.com/YahyaMudallal/newsWebSite/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ArticleService provides business logic for articles.
type ArticleService struct {
	articleRepo repositories.ArticleRepository
	userRepo repositories.UserRepository
	commentRepo repositories.CommentRepository
	newsClient clients.NewsClient
	geminiClient clients.GeminiClient
}

// NewArticleService creates a new ArticleService.
func NewArticleService(
	articleRepo repositories.ArticleRepository,
	userRepo repositories.UserRepository,
	commentRepo repositories.CommentRepository,
	newsClient clients.NewsClient,
	geminiClient clients.GeminiClient,
	) *ArticleService {
	return &ArticleService{
		articleRepo: articleRepo,
		userRepo: userRepo,
		commentRepo: commentRepo,
		newsClient: newsClient,
		geminiClient: geminiClient,
	}
}

// GetAllArticles retrieves all articles.
func (s *ArticleService) GetAllArticles(ctx context.Context) ([]models.Article, error) {
	return s.articleRepo.GetAll(ctx)
}

// GetArticleByID retrieves an article by its ID.
func (s *ArticleService) GetArticleByID(ctx context.Context, id bson.ObjectID) (*models.Article, error) {
	return s.articleRepo.GetByID(ctx, id)
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

	return s.articleRepo.Create(ctx, article)
}

// DeleteArticle deletes an article by its ID.
func (s *ArticleService) DeleteArticle(ctx context.Context, articleID bson.ObjectID, userID bson.ObjectID) error {

	// check if the article exists
	_, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("%w : article not found", apperrors.ErrNotFound)
	}

	// check if the user is admin
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w : user not found", apperrors.ErrNotFound)
	}
	if !user.IsAdmin {
		return fmt.Errorf("%w : unauthorized to delete this article", apperrors.ErrUnauthorized)
	}

	// delete the comments associated with the article
	err = s.commentRepo.DeleteByArticleID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("%w : failed to delete associated comments", apperrors.ErrInternal)
	}

	// delete the article
	return s.articleRepo.Delete(ctx, articleID)
}

// SyncDailyArticles Called daily by the server to Synchronize the articles database
// by calling the new api to get new articles and add them to the database.
func (s *ArticleService) SyncDailyArticles(ctx context.Context) error {
	// fetch the article from the client
	newArtciles, err := s.newsClient.FetchDailyArticles(ctx)
	if err != nil {
		return fmt.Errorf("%w : failed to fetch articles from news client", apperrors.ErrInternal)
	}

	// save the articles in the database (maybe add CreateMany in the repository for better performance)
	_, err = s.articleRepo.CreateMany(ctx, newArtciles)
	if err != nil {
		return fmt.Errorf("%w : failed to save articles in the database", apperrors.ErrInternal)
	}

	return nil
}

// GenerateSummary call the Gemini API to generate a summary (tldr) for the given article and update the article in the database with the new summary.
func (s *ArticleService) GenerateSummary(ctx context.Context, articleID bson.ObjectID) error {

	// get the article from the database
	article, err := s.articleRepo.GetByID(ctx, articleID)
	if err != nil {
		return fmt.Errorf("%w : failed to get the article from the database: %w", apperrors.ErrInternal, err)
	}

	// check if the article has already a summary
	if article.Summary != "" {
		return fmt.Errorf("%w : the article already has a summary", apperrors.ErrValidation)
	}

	// call the Gemini API to generate a summary for the article
	summary, err := s.geminiClient.GenerateTLDR(ctx, article)
	if err != nil {
		return fmt.Errorf("%w : failed to generate the summary: %w", apperrors.ErrInternal, err)
	}

	// add the summary to the article
	article.Summary = summary

	// update the article in the database
	err = s.articleRepo.Update(ctx, article)
	if err != nil {
		return fmt.Errorf("%w : failed to update the article in the database: %w", apperrors.ErrInternal, err)
	}

	return nil
}