package sub

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var (
// 	// cfgFile string
// 	rootCmd = &cobra.Command{
// 		Use:   "xfsmiddle",
// 		Short: "Git is a distributed version control system.",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			return cmd.Help()
// 			// Error(cmd, args, errors.New("unrecognized command"))
// 		},
// 	}
// )

// func Execute() {
// 	// if err := rootCmd.Execute(); err != nil {
// 	// 	_, _ = fmt.Fprintln(os.Stderr, err)
// 	// 	os.Exit(1)
// 	// }
// 	rootCmd.Execute()
// }
var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
