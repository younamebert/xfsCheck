package utils

import (
	"xfsmiddle"
	"xfsmiddle/db"
	"xfsmiddle/server/web"
)

type GroupsConfig []map[string][]string

type RpcServerConfig struct {
	Apihost string
	Timeout string
}

type RpcClientConfig struct {
	Apihost string
	Timeout string
}

type TokenDbConfig struct {
	Db db.IDatabase
}

func StartBack(g GroupsConfig, token TokenDbConfig, rs RpcServerConfig) error {
	server := xfsmiddle.NewRpcServer()

	groups := setupGroups(g)
	webtoken := setupToken(token, groups)
	server.RegisterName("Token", webtoken)
	if err := server.Start(rs.Apihost, rs.Timeout); err != nil {
		return err
	}
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
