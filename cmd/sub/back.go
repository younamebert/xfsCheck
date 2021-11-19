package sub

import (
	"xfsmiddle/db"
	"xfsmiddle/utils"
	util "xfsmiddle/utils"

	"github.com/spf13/cobra"
)

var (
	cfgFile          string
	rpcServerCommand = &cobra.Command{
		Use:   "start",
		Short: "start rpc server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Start()
		},
	}
	// rpcServerStartCommand = &cobra.Command{
	// 	Use:   "start",
	// 	Short: "rpc server start",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return tokenList()
	// 	},
	// }
	// rpcServerStopCommand = &cobra.Command{
	// 	Use:   "stop",
	// 	Short: "rpc server stop",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return newToken()
	// 	},
	// }
)

func Start() error {
	config, err := runConfig(cfgFile)
	if err != nil {
		return err
	}
	// utils.
	gConfig := setupGroupConfig(*config)

	db, err := db.New(config.tokenDb.tokenDbDir)
	if err != nil {
		return err
	}
	tConfig := util.TokenDbConfig{
		Db: db,
	}

	rpcsever := utils.RpcServerConfig{
		Apihost: config.serve.rpcServerApiHost,
		Timeout: config.serve.rpcServerApiTimeOut,
	}

	if err := util.StartBack(gConfig, tConfig, rpcsever); err != nil {
		return err
	}
	return nil
}

func setupGroupConfig(g backConfig) util.GroupsConfig {
	rights := make(util.GroupsConfig, 0)

	rights = append(rights, map[string][]string{
		"Chain": g.group.chain,
	})
	rights = append(rights, map[string][]string{
		"Miner": g.group.miner,
	})
	rights = append(rights, map[string][]string{
		"Net": g.group.net,
	})
	rights = append(rights, map[string][]string{
		"TxPool": g.group.txpool,
	})
	rights = append(rights, map[string][]string{
		"Wallet": g.group.wallet,
	})
	rights = append(rights, map[string][]string{
		"State": g.group.state,
	})

	return rights
}

func init() {
	// rpcServerCommand.AddCommand(rpcServerStopCommand)
	// rpcServerCommand.AddCommand(rpcServerStartCommand)
	rootCmd.AddCommand(rpcServerCommand)
	mFlags := rootCmd.PersistentFlags()
	mFlags.StringVarP(&cfgFile, "config", "C", "", "Set config file")
}
