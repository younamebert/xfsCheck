package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// cfgFile string
	rootCmd = &cobra.Command{
		Use: "xfsmiddle [flags] command [flags...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
