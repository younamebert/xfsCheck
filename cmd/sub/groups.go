package sub

import (
	"fmt"
	"xfsmiddle"
	"xfsmiddle/common"

	"github.com/spf13/cobra"
)

var (
	groupsCommand = &cobra.Command{
		Use:   "group",
		Short: "permission group",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	getGroupsCommand = &cobra.Command{
		Use:   "getgroups",
		Short: "get permission group configuration",
		RunE:  getGroups,
	}
)

func getGroups(cmd *cobra.Command, args []string) error {
	config, err := rpcClientConfigParams(cfgFile)
	if err != nil {
		return err
	}
	req := new(empty)
	cli := xfsmiddle.NewClient(config.rpcClientApiHost, config.rpcClientApiTimeOut)

	result := make([]map[string]interface{}, 0)
	if err := cli.Call("Groups.GetGroups", &req, &result); err != nil {
		return err
	}
	bs, err := common.MarshalIndent(result)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", string(bs))
	return nil
}

func init() {
	groupsCommand.AddCommand(getGroupsCommand)
	rootCmd.AddCommand(groupsCommand)
}
