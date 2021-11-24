package xfsmiddle

import (
	"context"
	"flag"
	"strings"
	"time"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
)

type Client struct {
	Apihost string
	Timeout string
}

func NewClient(clientHost, timeout string) *Client {
	return &Client{
		Apihost: clientHost,
		Timeout: timeout,
	}
}

func (cli *Client) Call(methods string, args, reply interface{}) error {

	temp := strings.Split(methods, ".")

	var addr = flag.String("addr", cli.Apihost, "client address")
	flag.Parse()
	d, err := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	if err != nil {
		return err
	}
	opt := client.DefaultOption
	opt.SerializeType = protocol.JSON

	timeDur, err := time.ParseDuration(cli.Timeout)
	if err != nil {
		return err
	}
	opt.ConnectTimeout = timeDur

	xclient := client.NewXClient(temp[0], client.Failtry, client.RandomSelect, d, opt)
	defer xclient.Close()
	err = xclient.Call(context.Background(), temp[1], args, reply)
	if err != nil {
		return err
	}
	return nil
}
