package clients

import (
	"context"

	"github.com/YahyaMudallal/newsWebSite/internal/models"
)

type GeminiClient interface {
	GenerateTLDR(ctx context.Context, article *models.Article) (string, error)
}

type GeminiAPIClient struct {
	apiKey string
}

func NewGeminiAPIClient(apiKey string) GeminiClient {
	return &GeminiAPIClient{apiKey: apiKey}
}

// GenerateTLDR generates a summary (tldr) for the given article using the Gemini API.
func (c *GeminiAPIClient) GenerateTLDR(ctx context.Context, article *models.Article) (string, error) {
	// TODO
	return "", nil
}