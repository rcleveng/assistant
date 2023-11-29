package llm

import (
	"context"
)

// LLM Client, currently only PALM is supported.
type LlmClient interface {
	GenerateText(ctx context.Context, prompt string) (string, error)
	EmbedText(ctx context.Context, text string) ([]float32, error)
	BatchEmbedText(ctx context.Context, text []string) ([][]float32, error)
	Close() error
}
