package embedcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embed is a commandline to add embeddedings for text",
	Long:  `Command to generate text embeddings`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("need command")
	},
}

var UseDatabase bool

func init() {
	RootCmd.PersistentFlags().BoolVarP(&UseDatabase, "use_database", "d", false, "Use the database")
}
