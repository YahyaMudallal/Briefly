package clients

import (
	"context"
	"fmt"

	"github.com/YahyaMudallal/newsWebSite/internal/apperrors"
	"github.com/YahyaMudallal/newsWebSite/internal/models"
	"google.golang.org/genai"
)

type GeminiClient interface {
	GenerateTLDR(ctx context.Context, article *models.Article) (string, error)
}

type GeminiAPIClient struct {
	apiKey string
	client *genai.Client
}

func NewGeminiAPIClient(ctx context.Context, apiKey string) (GeminiClient , error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("%w failed to create Gemini API client: %w", apperrors.ErrInternal, err)
	}
	return &GeminiAPIClient{apiKey: apiKey, client: client}, nil
}

// GenerateTLDR generates a summary (tldr) for the given article using the Gemini API.
func (c *GeminiAPIClient) GenerateTLDR(ctx context.Context, article *models.Article) (string, error) {

	// use the 1.5 flash model for speed and lower cost
	modelID := "gemini-flash-latest"

	// prompt engineering to generate a concise summary of the article
	prompt := fmt.Sprintf(`
        You are an expert news editor. Your task is to provide a concise TL;DR (summary) of the following article.
        
        Constraints:
        - Maximum 3 concise sentences.
        - Maintain a factual and neutral tone.
        - Focus on the key takeaway.

        Article Title: %s
        Article Description: %s
        Article Content: %s
    `, article.Title, article.Description, article.Content)

	temp := float32(0.3) // low temperature for more factual output
	config := &genai.GenerateContentConfig{
		Temperature: &temp,
	}

	// call the Gemini API to generate the summary
	resp, err := c.client.Models.GenerateContent(ctx, modelID, genai.Text(prompt), config)
	if err != nil {
		return "", fmt.Errorf("%w gemini generation failed: %w", apperrors.ErrInternal, err)
	}
	
	// extract the generated summary from the response
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return resp.Candidates[0].Content.Parts[0].Text, nil
	}

	// if we reach here, it means the response did not contain the expected content
	return "", fmt.Errorf("%w gemini generation failed: no content in response", apperrors.ErrInternal)
}
