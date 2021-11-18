package web

import (
	"context"
	"fmt"
	"xfsmiddle"
)

type Token struct {
	TokenManage *xfsmiddle.TokenManage
	Rights      *xfsmiddle.Groups
}

type DelTokenArgs struct {
	Token string `json:"token"`
}

func (t *Token) DelToken(ctx context.Context, args *DelTokenArgs, reply *string) error {
	fmt.Println(args.Token)
	*reply = args.Token
	return nil
}

func (t *Token) NewToken(ctx context.Context, args *DelTokenArgs, reply *string) error {
	return nil
}
