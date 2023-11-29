package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rcleveng/assistant/server/env"
	"github.com/rcleveng/assistant/server/llm"
	"github.com/tmc/langchaingo/textsplitter"
)

func main() {
	env, err := env.NewServerEnvironment(env.COMMANDLINE)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
		os.Exit(2)
	}

	ctx := context.Background()
	client, err := llm.NewPalmLLMClient(ctx, env)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
		os.Exit(2)
	}
	defer client.Close()

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkOverlap = 0
	splitter.ChunkSize = 20

	splits, err := splitter.SplitText("What is the world's largest island that's not a continent?")
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
		os.Exit(2)
	}

	for _, split := range splits {

		resp, err := client.EmbedText(ctx, split)
		if err != nil {
			fmt.Printf("ERROR: %v\n\n", err)
			continue
		}

		fmt.Printf("Embedding: [%#v] '%s]\n", resp, split)
	}

}
