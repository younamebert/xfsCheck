package web

import (
	"context"
	"errors"
	"strings"
	"xfsmiddle"
)

// var tokenErrorCode = -10086

type Token struct {
	TokenManage *xfsmiddle.TokenManage
	Rights      *xfsmiddle.Groups
}

type DelTokenArgs struct {
	Token string `json:"token"`
}

func (t *Token) DelToken(ctx context.Context, args *DelTokenArgs, reply *string) error {
	if args.Token == "" {
		return ctx.Err()
	}
	return t.TokenManage.DelToken(args.Token)
}

func (t *Token) NewToken(ctx context.Context, args *NewTokenArgs, reply *string) error {
	var gs string

	gs = strings.Join(t.Rights.GetTypes(), ",")
	if args.Group != "" {
		gs = args.Group
		want := strings.Split(args.Group, ",")
		got := t.Rights.Get(want)
		if len(want) != len(got) {
			return errors.New("expected value is not reached")
		}
	}

	result, err := t.TokenManage.NewToken(gs)
	if err != nil {
		return err
	}

	*reply = result
	return nil
}

func (t *Token) GetGroupByToken(ctx context.Context, args *GetGroupByTokenArgs, reply *string) error {
	if args.Token == "" {
		return ctx.Err()
	}
	gouroup, err := t.TokenManage.GetToken(args.Token)
	if err != nil {
		return err
	}
	*reply = string(gouroup)
	return nil
}

func (t *Token) PutTokenGroup(ctx context.Context, args *PutTokenGroupArgs, reply *string) error {

	if args.Group != "" {
		want := strings.Split(args.Group, ",")
		got := t.Rights.Get(want)
		if len(want) != len(got) {
			return errors.New("expected value is not reached")
		}
	}
	if _, err := t.TokenManage.GetToken(args.Token); err != nil {
		return err
	}
	return t.TokenManage.SetTokenGroup(args.Token, args.Group)
}
