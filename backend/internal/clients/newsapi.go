package clients

import (
	"context"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
)

type NewsClient interface {
    FetchDailyArticles(ctx context.Context) ([]models.Article, error)
}	

type newsAPIClient struct {
    apiKey string
}

func NewNewsAPIClient(apiKey string) NewsClient {
    return &newsAPIClient{apiKey: apiKey}
}

func (c *newsAPIClient) FetchDailyArticles(ctx context.Context) ([]models.Article, error) {
	// TODO
	return nil, nil
}