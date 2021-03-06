package main

import (
	"fmt"
	"net"

	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

type GenericClient interface {
	Call(method string, arg interface{}, res interface{}) error
}

//---------------------------------------------------------------------

type ArithClient struct {
	GenericClient
}

func (a ArithClient) Add(arg AddArgs) (ret int, err error) {
	err = a.Call("test.1.arith.add", arg, &ret)
	return
}

func (a ArithClient) Broken() (err error) {
	err = a.Call("test.1.arith.broken", nil, nil)
	return
}

//---------------------------------------------------------------------

type Client struct {
	port int
}

func (s *Client) Run() (err error) {
	var c net.Conn
	if c, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", s.port)); err != nil {
		return
	}

	xp := rpc2.NewTransport(c, nil, nil)
	cli := ArithClient{GenericClient: rpc2.NewClient(xp, nil)}

	for A := 10; A < 23; A += 2 {
		var res int
		if res, err = cli.Add(AddArgs{A: A, B: 34}); err != nil {
			return
		}
		fmt.Printf("result is -> %v\n", res)
	}

	err = cli.Broken()
	fmt.Printf("for broken: %v\n", err)

	return nil
}
