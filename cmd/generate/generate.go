package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rcleveng/assistant/server/env"
	"github.com/rcleveng/assistant/server/llm"
)

func main() {
	args := os.Args[1:]

	env := &env.ServerEnvironment{
		PalmApiKey:           os.Getenv("PALM_KEY"),
		DatabaseConnection:   "",
		DatabaseUserName:     "",
		DatabasePassword:     "",
		ExecutionEnvironment: env.COMMANDLINE,
	}

	ctx := context.Background()
	llmclient, err := llm.NewPalmLLMClient(ctx, env)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
		os.Exit(2)
	}

	for _, arg := range args {
		response, err := llmclient.GenerateText(ctx, arg)
		if err != nil {
			fmt.Printf("ERROR: %v\n\n", err)
			continue
		}
		fmt.Println(response)
	}

}