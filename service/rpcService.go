package service

import (
	"fmt"
	commonChannel "k8sOp/channel"
	"net"
	"net/rpc"
)

type rpcService struct{}

type RequestRpc struct{}

type Args struct {
	ServiceName string
}

type Reply struct {
	Msg string
}

var RpcService rpcService

func (c *rpcService) Start() {
	fmt.Println("rpc server running at 5556 port")

	rpc.Register(new(RequestRpc))

	ln, err := net.Listen("tcp", ":5556")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}
}

func (c *RequestRpc) AddRequest(args *Args, reply *Reply) error {
	fmt.Printf("in service :: %s.\n", args.ServiceName)
	commonChannel.CommonChannel.RequestChannel <- args.ServiceName
	reply.Msg = "OK"
	return nil
}
