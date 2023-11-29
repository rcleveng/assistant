package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/rcleveng/assistant/server/env"
)

func createResponse(statusCode int, resp any) (func(req *http.Request) (*http.Response, error), error) {
	body, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	doer := func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			Status:        "",
			StatusCode:    statusCode,
			Header:        map[string][]string{},
			Body:          io.NopCloser(bytes.NewReader(body)),
			ContentLength: int64(len(body)),
			Request:       req,
		}
		return resp, nil
	}
	return doer, nil
}

func TestGenerateText(t *testing.T) {
	message := "Hello World"
	doer, err := createResponse(200, &GenerateTextResponse{
		Candidates: []*TextCompletion{
			{
				Output: message,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}

	client := &PalmLLMClient{
		environment: &env.ServerEnvironment{
			PalmApiKey:           "",
			DatabaseConnection:   "",
			DatabaseUserName:     "",
			DatabasePassword:     "",
			ExecutionEnvironment: env.GOTEST,
		},
		endpoint: "https://generativelanguage.googleapis.com/v1beta3",
		doer:     DoerFunc(doer),
	}

	resp, err := client.GenerateText(context.Background(), "test")
	if err != nil {
		t.Error(err)
	}

	if string(resp) != message {
		t.Errorf("Expected '%s' got '%s'", message, string(resp))
	}
}
