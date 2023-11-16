package main

import (
	"context"
	"fmt"
	"os"

	gl "cloud.google.com/go/ai/generativelanguage/apiv1beta2"
	pb "cloud.google.com/go/ai/generativelanguage/apiv1beta2/generativelanguagepb"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := gl.NewTextRESTClient(ctx, option.WithAPIKey(os.Getenv("PALM_KEY")))
	if err != nil {
		panic(err)
	}

	defer client.Close()

	req := &pb.EmbedTextRequest{
		Model: "models/embedding-gecko-001",
		Text:  "What is the world's largest island that's not a continent?",
	}

	resp, err := client.EmbedText(ctx, req)
	if err != nil {
		panic(err)
	}

	e := resp.GetEmbedding().GetValue()
	fmt.Printf("Embedding: %#v\n", e)
}
