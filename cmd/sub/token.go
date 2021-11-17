package sub

import "github.com/spf13/cobra"

var (
	tokenCommand = &cobra.Command{
		Use:   "token",
		Short: "get token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	tokenListCommand = &cobra.Command{
		Use:   "list",
		Short: "get token list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return tokenList()
		},
	}
	tokenNewCommand = &cobra.Command{
		Use:   "new",
		Short: "Create token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newToken()
		},
	}
	tokenDelCommand = &cobra.Command{
		Use:   "delete <token>",
		Short: "Delete token",
		RunE:  delToken,
	}
)

func tokenList() error {
	return nil
}

func newToken() error {
	return nil
}

func delToken(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	tokenCommand.AddCommand(tokenListCommand)
	tokenCommand.AddCommand(tokenDelCommand)
	tokenCommand.AddCommand(tokenNewCommand)
	rootCmd.AddCommand(tokenCommand)
}
