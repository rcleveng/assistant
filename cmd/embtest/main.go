package main

import (
	"context"
	"fmt"
	"os"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"github.com/tmc/langchaingo/textsplitter"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := gl.NewTextRESTClient(ctx, option.WithAPIKey(os.Getenv("PALM_KEY")))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkOverlap = 0
	splitter.ChunkSize = 20

	splits, err := splitter.SplitText("What is the world's largest island that's not a continent?")
	if err != nil {
		panic(err)
	}

	for _, split := range splits {

		req := &pb.EmbedTextRequest{
			Model: "models/embedding-gecko-001",
			Text:  split,
		}

		resp, err := client.EmbedText(ctx, req)
		if err != nil {
			panic(err)
		}

		e := resp.GetEmbedding().GetValue()
		fmt.Printf("Embedding: [%#v] '%s]\n", e, split)
	}

}
