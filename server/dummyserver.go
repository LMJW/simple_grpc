package main

import (
	context "context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type dummyServer struct{}

func (s *dummyServer) GetDummy(ctx context.Context, in *DummyRequest) (*DummyResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("Unauthenticated")
	}

	fmt.Println(p.Addr.Network(), p.Addr.String())

	auth, ok := p.AuthInfo.(credentials.TLSInfo)

	fmt.Println(auth.State.ServerName)
	fmt.Printf("%s\n", auth.State.TLSUnique)

	for i, cert := range auth.State.PeerCertificates {
		fmt.Println(i)
		fmt.Printf("Peer certificate %v, Issued by %v", cert.Subject.CommonName, cert.Issuer.CommonName)
	}

	for i, chain := range auth.State.VerifiedChains {
		fmt.Println(i, "!!!!!!!!!!!!!!!!!!!!!!!!!!")
		for j, cert := range chain {
			fmt.Println(j)
			fmt.Printf("Chain %v::::Peer certificate %v, Issued by %v", j, cert.Subject.CommonName, cert.Issuer.CommonName)
		}
	}
	fmt.Println("?")
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
