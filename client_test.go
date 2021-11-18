package xfsmiddle

import (
	"fmt"
	"testing"
)

type testDelTokenArgs struct {
	Token string `json:"token"`
}

func Test_Client(t *testing.T) {

	cli := NewClient("localhost:9002", "60s")
	args := testDelTokenArgs{
		Token: "2112122121",
	}
	reply := new(string)
	if err := cli.Call("Token.DelToken", args, reply); err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(*reply)
}
