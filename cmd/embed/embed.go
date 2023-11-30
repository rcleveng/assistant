package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
	"github.com/rcleveng/assistant/server/db"
	"github.com/rcleveng/assistant/server/env"
	"github.com/rcleveng/assistant/server/llm"
	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/textsplitter"
)

var rootCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embed is a commandline to add embeddedings for text",
	Long:  `Command to generate text embeddings`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := env.NewEnvironment(env.COMMANDLINE)
		if err != nil {
			return err
		}

		if err := embed(env, args); err != nil {
			return err
		}

		return nil
	},
}

var UseDatabase bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&UseDatabase, "use_database", "d", false, "Use the database")
}

type Splitter interface {
	SplitText(text string) ([]string, error)
}

func embedAndAdd(ctx context.Context, splitter Splitter, lm llm.LlmClient, db db.EmbeddingsDB, texts []string) error {
	splits := make([]string, 0, len(texts))
	for _, text := range texts {
		cursplits, err := splitter.SplitText(text)
		if err != nil {
			return err
		}
		splits = append(splits, cursplits...)
	}

	resp, err := lm.BatchEmbedText(ctx, splits)
	if err != nil {
		fmt.Printf("ERROR: %v\n\n", err)
	}

	glog.V(2).Infof("%s", spew.Sdump(resp))

	for i, e := range resp {
		author := int64(0)
		if i < len(splits) {
			glog.V(1).Infof("Embedding: [%d] [%#v] '%s']\n", i, e, splits[i])
			if _, err := db.Add(author, splits[i], e); err != nil {
				return err
			}
		} else {
			fmt.Printf("past end of splits with i=%d\n", i)
			fmt.Printf("Past Split Embedding: [%d] [%#v]]\n", i, e)
		}
	}
	return nil
}

func embed(env *env.Environment, text []string) error {
	ctx := context.Background()
	llm, err := llm.NewPalmLLMClient(ctx, env)
	if err != nil {
		return err
	}
	defer llm.Close()

	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkOverlap = 0
	splitter.ChunkSize = 20

	var edb db.EmbeddingsDB

	if UseDatabase {
		fmt.Println("Using postgresql database")
		edb, err = db.NewPostgresDatabase(env)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Skipping database")
		edb = NoopEmbeddingsDB{}
	}
	defer edb.Close()

	if err = embedAndAdd(ctx, splitter, llm, edb, text); err != nil {
		return err
	}

	return nil
}

type NoopEmbeddingsDB struct{}

func (n NoopEmbeddingsDB) Close() {}
func (n NoopEmbeddingsDB) Add(author int64, text string, embeddings []float32) (int64, error) {
	return 0, nil
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
