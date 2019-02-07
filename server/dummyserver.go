package main

import (
	context "context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type dummyServer struct{}

func (s *dummyServer) GetDummy(ctx context.Context, in *DummyRequest) (*DummyResponse, error) {
	fmt.Println(in)
	return &DummyResponse{Msg: "received"}, nil
}

func main() {
	creds, _ := credentials.NewServerTLSFromFile("./server.crt", "./server.key")
	s, err := net.Listen("tcp", ":54332")
	if err != nil {
		log.Fatal(err)
	}
	gs := grpc.NewServer(grpc.Creds(creds))
	RegisterDummyServiceServer(gs, &dummyServer{})

	gs.Serve(s)
}
