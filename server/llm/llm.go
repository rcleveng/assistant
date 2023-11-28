package llm

import (
	"context"
)

type LlmClient interface {
	GenerateText(ctx context.Context, prompt string) (string, error)
	EmbedText(ctx context.Context, text string) ([]float32, error)
	BatchEmbedText(ctx context.Context, text []string) ([][]float32, error)
	Close() error
}
