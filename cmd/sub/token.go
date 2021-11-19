package sub

import (
	"xfsmiddle"

	"github.com/spf13/cobra"
)

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
		Use:   "new [rights_group]",
		Short: "Create token",
		RunE:  newToken,
	}
	tokenDelCommand = &cobra.Command{
		Use:   "delete <token>",
		Short: "Delete token",
		RunE:  delToken,
	}
)

func tokenList() error {
	// config, err := rpcClientConfigParams(cfgFile)
	// if err != nil {
	// 	return err
	// }
	// cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)
	// cli.Call()
	return nil
}

func newToken(cmd *cobra.Command, args []string) error {
	config, err := rpcClientConfigParams(cfgFile)
	if err != nil {
		return err
	}

	cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)

	req := new(newTokenArgs)
	if len(args) > 0 {
		req.group = args[0]
	}

	var result string
	if err := cli.Call("Token.NewToken", &req, &result); err != nil {
		return err
	}
	return nil
}

func delToken(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return cmd.Help()
	}
	config, err := rpcClientConfigParams(cfgFile)
	if err != nil {
		return err
	}

	cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)

	req := &delTokenArgs{
		token: args[0],
	}

	var result string
	if err := cli.Call("Token.DelToken", &req, &result); err != nil {
		return err
	}
	return nil
}

func init() {
	tokenCommand.AddCommand(tokenListCommand)
	tokenCommand.AddCommand(tokenDelCommand)
	tokenCommand.AddCommand(tokenNewCommand)
	rootCmd.AddCommand(tokenCommand)
}
