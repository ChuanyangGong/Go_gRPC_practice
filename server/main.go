package main

import (
	"go_grpc_practice/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	g := grpc.NewServer()
	proto.RegisterLoginerServer(g, newServer())
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	err = g.Serve(listener)
	if err != nil {
		panic(err)
	}
}
