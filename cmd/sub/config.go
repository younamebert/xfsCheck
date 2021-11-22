package sub

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "./config.yml"
)

type backConfig struct {
	serve   serverConfig
	group   groupConfig
	tokenDb tokenDbConfig
	gateway gatewayConfig
}

type serverConfig struct {
	rpcServerApiHost    string
	rpcServerApiTimeOut string
}

type clientConfig struct {
	rpcClientApiHost    string
	rpcClientApiTimeOut string
}

type groupConfig struct {
	miner  []string
	wallet []string
	chain  []string
	state  []string
	txpool []string
	net    []string
}

type tokenDbConfig struct {
	tokenDbDir string
}

type gatewayConfig struct {
	apihost  string
	timeout  string
	nodeaddr string
	rpcaddr  string
}

func runConfig(configPath string) (*backConfig, error) {
	config := viper.New()
	if err := setViperPath(config, configPath); err != nil {
		return nil, err
	}
	return &backConfig{
		serve:   rpcServerConfigParams(config),
		group:   groupConfigParams(config),
		tokenDb: tokenDbConfigParams(config),
		gateway: gatewayConfigParams(config),
	}, nil

}
func setViperPath(v *viper.Viper, customFile string) error {
	filename := filepath.Base(defaultConfigFile)
	ext := filepath.Ext(defaultConfigFile)
	configPath := filepath.Dir(defaultConfigFile)
	v.AddConfigPath(configPath)
	v.SetConfigType(strings.TrimPrefix(ext, "."))
	v.SetConfigName(strings.TrimSuffix(filename, ext))
	v.SetConfigFile(customFile)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func rpcServerConfigParams(v *viper.Viper) serverConfig {
	serverConfig := serverConfig{}
	serverConfig.rpcServerApiHost = v.GetString("rpcserver.rpcserver")
	serverConfig.rpcServerApiTimeOut = v.GetString("rpcserver.timeout")
	return serverConfig
}

func rpcClientConfigParams(configPath string) (*clientConfig, error) {
	config := viper.New()
	if err := setViperPath(config, configPath); err != nil {
		return nil, err
	}

	return &clientConfig{
		rpcClientApiHost:    config.GetString("rpcclient.rpcclient"),
		rpcClientApiTimeOut: config.GetString("rpcclient.timeout"),
	}, nil
}

func groupConfigParams(v *viper.Viper) groupConfig {
	groupConfig := groupConfig{}
	groupConfig.chain = v.GetStringSlice("group.chain")
	groupConfig.miner = v.GetStringSlice("group.miner")
	groupConfig.state = v.GetStringSlice("group.state")
	groupConfig.net = v.GetStringSlice("group.net")
	groupConfig.txpool = v.GetStringSlice("group.txpool")
	groupConfig.wallet = v.GetStringSlice("group.wallet")

	return groupConfig
}

func tokenDbConfigParams(v *viper.Viper) tokenDbConfig {
	tokenDbConfig := tokenDbConfig{}
	datadbDir := v.GetString("leveldb.datadir")
	tokenDbConfig.tokenDbDir = filepath.Join(datadbDir, v.GetString("leveldb.tokendir"))
	return tokenDbConfig
}

func gatewayConfigParams(v *viper.Viper) gatewayConfig {
	gatewayConfig := gatewayConfig{}
	gatewayConfig.apihost = v.GetString("gateway.apihost")
	gatewayConfig.timeout = v.GetString("gateway.timeout")
	gatewayConfig.rpcaddr = v.GetString("gateway.rpcaddr")
	gatewayConfig.nodeaddr = v.GetString("gateway.nodeaddr")
	return gatewayConfig
}
