package web

import (
	"context"
	"xfsmiddle"
)

type jsonRPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResp struct {
	JSONRPC string              `json:"jsonrpc"`
	Result  interface{}         `json:"result"`
	Error   *xfsmiddle.RPCError `json:"error"`
	ID      int                 `json:"id"`
}

type ForWard struct {
	TokenManage *xfsmiddle.TokenManage
}

func SendForWard(ctx context.Context, args *jsonRPCReq, reply *string) error {
	return nil
}
