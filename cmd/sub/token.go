package sub

import (
	"fmt"
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
	tokenNewCommand = &cobra.Command{
		Use:   "new [rights_group]",
		Short: "create permission token",
		RunE:  newToken,
	}
	tokenDelCommand = &cobra.Command{
		Use:   "delete <token>",
		Short: "delete the specified token",
		RunE:  delToken,
	}
	tokenUpdateCommand = &cobra.Command{
		Use:   "settoken <token> <group>",
		Short: "modify token permissions",
		RunE:  setToken,
	}
	getTokenGroupCommand = &cobra.Command{
		Use:   "gettokengroup <token>",
		Short: "get the specified token permission",
		RunE:  getGroupByToken,
	}
)

func getGroupByToken(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}
	config, err := rpcClientConfigParams(cfgFile)
	if err != nil {
		return err
	}
	req := new(getTokenGroupArgs)
	req.token = args[0]

	cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)

	var result string
	if err := cli.Call("Token.GetGroupByToken", &req, &result); err != nil {
		return err
	}
	fmt.Printf("%v\n", result)
	return nil
}

func setToken(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmd.Help()
	}
	config, err := rpcClientConfigParams(cfgFile)
	if err != nil {
		return err
	}
	req := new(setTokenArgs)
	req.group = args[0]
	req.token = args[1]
	cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)

	var result string
	if err := cli.Call("Token.PutTokenGroup", &req, &result); err != nil {
		return err
	}
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
	fmt.Printf("%s\n", result)
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
	tokenCommand.AddCommand(tokenDelCommand)
	tokenCommand.AddCommand(tokenNewCommand)
	tokenCommand.AddCommand(tokenUpdateCommand)
	tokenCommand.AddCommand(getTokenGroupCommand)
	rootCmd.AddCommand(tokenCommand)
}
