package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rcleveng/assistant/server/env"
)

// Single Method Interfacr helper.
// https://eli.thegreenplace.net/2023/the-power-of-single-method-interfaces-in-go/
type DoerFunc func(*http.Request) (*http.Response, error)

func (d DoerFunc) Do(r *http.Request) (*http.Response, error) {
	return d(r)
}

// Mock helper for HTTP requests (all call Do on the client)
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type PalmLLMClient struct {
	// Everything we need context wise from the environment
	environment *env.Environment
	// base url endpoint for this service
	endpoint string
	// Used to mock out HttpClient.Do
	doer Doer
}

func (c *PalmLLMClient) Close() error {
	return nil
}

func (c *PalmLLMClient) Do(httpMethod, model, fn string, request, response any) error {
	url := fmt.Sprintf("%s/%s:%s", c.endpoint, model, fn)

	body, err := json.Marshal(request)
	if err != nil {
		return err
	}

	var httpReq *http.Request = nil

	switch httpMethod {
	case http.MethodPost:
		httpReq, err = http.NewRequest(httpMethod, url, bytes.NewReader(body))
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unsupported method: %s", httpMethod)
	}
	httpReq.Header.Add("x-goog-api-key", c.environment.PalmApiKey)
	httpClient := c.doer
	if httpClient == nil {
		// Duck typing FTW
		httpClient = http.DefaultClient
	}
	res, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		// We have an error here. TODO- Parse error status
		return fmt.Errorf("returned HTTP error %d : %s", res.StatusCode, string(b))
	}

	if err = json.Unmarshal(b, &response); err != nil {
		return err
	}

	return nil
}

func (c *PalmLLMClient) Post(model, fn string, request, response any) error {
	return c.Do(http.MethodPost, model, fn, request, response)
}

func (c *PalmLLMClient) GenerateText(ctx context.Context, prompt string) (string, error) {

	req := &GenerateTextRequest{
		Prompt: &TextPrompt{
			Text: prompt,
		},
	}

	model := "models/text-bison-001"
	resp := &GenerateTextResponse{}

	if err := c.Post(model, "generateText", req, resp); err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 {
		s := resp.Candidates[0].Output
		return s, nil
	}
	return "", fmt.Errorf("no candidate response, just %#v", resp)
}

func (c *PalmLLMClient) EmbedText(ctx context.Context, text string) ([]float32, error) {
	model := "models/embedding-gecko-001"
	req := &EmbedTextRequest{
		Text: text,
	}

	resp := &EmbedTextResponse{}
	if err := c.Post(model, "embedText", req, resp); err != nil {
		return nil, err
	}

	if resp.Embedding == nil {
		return nil, fmt.Errorf("unable to find embedding json structure")
	}

	emb := resp.Embedding.Value
	return emb, nil
}

func (c *PalmLLMClient) BatchEmbedText(ctx context.Context, texts []string) ([][]float32, error) {
	model := "models/embedding-gecko-001"
	req := &BatchEmbedTextRequest{
		Texts: texts,
	}

	resp := &BatchEmbedTextResponse{}
	if err := c.Post(model, "batchEmbedText", req, resp); err != nil {
		return nil, err
	}

	if resp.Embeddings == nil {
		return nil, fmt.Errorf("unable to find embedding json structure")
	}

	embeddings := make([][]float32, 0, len(texts))
	for _, emb := range resp.Embeddings {
		embeddings = append(embeddings, emb.Value)
	}
	return embeddings, nil
}

func NewPalmLLMClient(ctx context.Context, environment *env.Environment) (*PalmLLMClient, error) {
	return &PalmLLMClient{
		environment: environment,
		endpoint:    "https://generativelanguage.googleapis.com/v1beta3",
	}, nil
}
