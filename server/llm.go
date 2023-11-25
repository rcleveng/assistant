package server

import (
	"context"
	"fmt"

	generativelanguage "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/rcleveng/assistant/server/env"
	"google.golang.org/api/option"
)

type LlmClient interface {
	Call(ctx context.Context, prompt string) (string, error)
	Close() error
}

type PalmLLMClient struct {
	c   *generativelanguage.TextClient
	ctx context.Context
}

func (c *PalmLLMClient) Close() error {
	return c.c.Close()
}

func (c *PalmLLMClient) Call(ctx context.Context, prompt string) (string, error) {
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
		fmt.Println("LLM Response: " + s)
		return s, nil
	}
	return "", fmt.Errorf("no candidate response, just %#v", resp)
}

func NewPalmLLMClient(ctx context.Context) (*PalmLLMClient, error) {
	env, ok := env.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unable to find serverenv on context")
	}
	palmKey, err := env.PalmApiKey()
	if err != nil {
		return nil, fmt.Errorf("error getting PALM API key")
	}

	apiKey := option.WithAPIKey(palmKey)
	c, err := generativelanguage.NewTextRESTClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return &PalmLLMClient{
		c:   c,
		ctx: ctx,
	}, nil
}

func OneShotSendToLLM(ctx context.Context, prompt string) (string, error) {
	c, err := NewPalmLLMClient(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()

	return c.Call(ctx, prompt)
}
