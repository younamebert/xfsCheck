package xfsmiddle

import (
	"flag"
	"xfsmiddle/logs"

	"github.com/smallnest/rpcx/server"
)

type Rpcserver struct {
	Serve *server.Server
	Logs  logs.ILogger
}

func NewRpcServer() *Rpcserver {
	return &Rpcserver{
		Serve: server.NewServer(),
		Logs:  logs.NewLogger("rpcserver"),
	}
}

func (s *Rpcserver) RegisterName(name string, rcvr interface{}) error {
	return s.Serve.RegisterName(name, rcvr, "")
}

func (s *Rpcserver) Start(Apihost, Timeout string) error {
	addr := flag.String("addr", Apihost, "server address")
	if err := s.Serve.Serve("tcp", *addr); err != nil {
		return err
	}
	return nil
}
