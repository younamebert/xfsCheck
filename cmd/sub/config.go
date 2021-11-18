package sub

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigFile = "./config.yml"
)

type Config struct {
	serve   serverConfig
	client  clientConfig
	group   groupConfig
	tokenDb tokenDbConfig
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

func runConfig(configPath string) (*Config, error) {
	config := viper.New()
	if err := setViperPath(config, configPath); err != nil {
		return nil, err
	}
	return &Config{
		serve:   rpcServerConfigParams(config),
		client:  rpcClientConfigParams(config),
		group:   groupConfigParams(config),
		tokenDb: tokenDbConfigParams(config),
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

func rpcClientConfigParams(v *viper.Viper) clientConfig {
	clientConfig := clientConfig{}
	clientConfig.rpcClientApiHost = v.GetString("rpcclient.rpcclient")
	clientConfig.rpcClientApiTimeOut = v.GetString("rpcclient.timeout")
	return clientConfig
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
