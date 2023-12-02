package embedcmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "remove an embedding",
	RunE: func(cmd *cobra.Command, args []string) error {
		// env, err := env.NewEnvironment(env.COMMANDLINE)
		// if err != nil {
		// 	return err
		// }

		return fmt.Errorf("implement me")
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
