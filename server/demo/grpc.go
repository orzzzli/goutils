package main

import (
	"context"
	"errors"
	"fmt"

	server2 "github.com/orzzzli/goutils/server"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (*server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.New("get error")
	}
	return &helloworld.HelloReply{
		Message: in.Name,
	}, nil
}

func main() {
	s := server2.NewServer(server2.ServeTypeGRPC, "test", "localhost:8888")
	err := s.RegisterService("testServer", (*helloworld.GreeterServer)(nil), &server{})
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
}
