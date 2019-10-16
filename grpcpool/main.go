package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/resolver"
)

func main() {
	go createServe()
	time.Sleep(time.Second)
	r, _ := GenerateAndRegisterManualResolver()
	r.InitialAddrs([]resolver.Address{resolver.Address{Addr: "127.0.0.1:12346", Metadata: int(10)}})
	conn, err := grpc.Dial(fmt.Sprintf("%s:///srv", r.Scheme()), grpc.WithBalancerName("my_round_robin"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	c := helloworld.NewGreeterClient(conn)

	req := &helloworld.HelloRequest{}
	c.SayHello(context.Background(), req)
}

type server struct{}

func (*server) SayHello(context.Context, *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{}, nil
}
func createServe() {
	l, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &server{})
	s.Serve(l)
}
