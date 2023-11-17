package server

import (
	"context"
	"fmt"
	"os"

	generativelanguage "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"google.golang.org/api/option"
)

type PalmLLMClient struct {
	c   *generativelanguage.TextClient
	ctx context.Context
}

func (c *PalmLLMClient) Close() error {
	return c.c.Close()
}

func (c *PalmLLMClient) Call(prompt string) (string, error) {
	req := &pb.GenerateTextRequest{
		Model: "models/text-bison-001",
		Prompt: &pb.TextPrompt{
			Text: prompt,
		},
	}

	resp, err := c.c.GenerateText(c.ctx, req)
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
	// This snippet has been automatically generated and should be regarded as a code template only.
	// It will require modifications to work:
	// - It may require correct/in-range values for request initialization.
	// - It may require specifying regional endpoints when creating the service client as shown in:
	//   https://pkg.go.dev/cloud.google.com/go#hdr-Client_Options
	apiKey := option.WithAPIKey(os.Getenv("PALM_KEY"))
	c, err := generativelanguage.NewTextRESTClient(ctx, apiKey)
	if err != nil {
		return nil, err
	}
	return &PalmLLMClient{
		c:   c,
		ctx: ctx,
	}, nil
}

func OneShotSendToLLM(prompt string) (string, error) {
	ctx := context.Background()
	c, err := NewPalmLLMClient(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()

	return c.Call(prompt)
}
