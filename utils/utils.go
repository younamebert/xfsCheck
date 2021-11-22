package utils

import (
	"xfsmiddle"
	"xfsmiddle/db"
	"xfsmiddle/server/web"
)

type GroupsConfig []map[string][]string

type RpcServerConfig struct {
	ApiHost string
	TimeOut string
}

type RpcClientConfig struct {
	ApiHost string
	TimeOut string
}

type TokenDbConfig struct {
	Db db.IDatabase
}

type GatewayConfig struct {
	RpcAddr  string
	TimeOut  string
	NodeAddr string
	ApiHost  string
}

func StartBack(group GroupsConfig, token TokenDbConfig, rpcserve RpcServerConfig, gates GatewayConfig) error {
	server := xfsmiddle.NewRpcServer()

	groups := setupGroups(group)
	webtoken := setupToken(token, groups)

	server.RegisterName("Token", webtoken)

	go func() {
		server.Start(rpcserve.ApiHost, rpcserve.TimeOut)
	}()

	gateway := xfsmiddle.NewRpcGateway(gates.RpcAddr, gates.ApiHost, gates.TimeOut, gates.NodeAddr, xfsmiddle.New(token.Db), server.ServiceMap())

	gateway.Start()
	select {}
	return nil
}

func setupGroups(g GroupsConfig) *xfsmiddle.Groups {
	groups := &xfsmiddle.Groups{}
	for _, v := range g {
		for ks, vs := range v {
			groups.Rights = append(groups.Rights, xfsmiddle.NewGroup(ks, vs))
		}
	}
	return groups
}

func setupToken(token TokenDbConfig, g *xfsmiddle.Groups) *web.Token {
	tokenHander := &web.Token{
		TokenManage: xfsmiddle.New(token.Db),
		Rights:      g,
	}
	return tokenHander
}
