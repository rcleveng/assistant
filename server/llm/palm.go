package llm

import (
	"context"
	"fmt"

	generativelanguage "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/rcleveng/assistant/server/env"
	"google.golang.org/api/option"
)

type PalmLLMClient struct {
	c           *generativelanguage.TextClient
	environment *env.ServerEnvironment
}

func (c *PalmLLMClient) Close() error {
	return c.c.Close()
}

func (c *PalmLLMClient) GenerateText(ctx context.Context, prompt string) (string, error) {
	req := &pb.GenerateTextRequest{
		Model: "models/text-bison-001",
		Prompt: &pb.TextPrompt{
			Text: prompt,
		},
	}

	resp, err := c.c.GenerateText(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 {
		s := resp.Candidates[0].Output
		//fmt.Println("LLM Response: " + s)
		return s, nil
	}
	return "", fmt.Errorf("no candidate response, just %#v", resp)
}

func (c *PalmLLMClient) EmbedText(ctx context.Context, text string) ([]float32, error) {
	return nil, nil
}

func (c *PalmLLMClient) BatchEmbedText(ctx context.Context, text []string) ([][]float32, error) {
	return nil, nil
}

func NewPalmLLMClient(ctx context.Context, environment *env.ServerEnvironment) (*PalmLLMClient, error) {
	apiKey := option.WithAPIKey(environment.PalmApiKey)
	c, err := generativelanguage.NewTextRESTClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return &PalmLLMClient{
		c:           c,
		environment: environment,
	}, nil
}
