package main

import (
	context "context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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
	log.Println("running")
	// creds, _ := credentials.NewServerTLSFromFile("./server.crt", "./server.key")

	server, err := tls.LoadX509KeyPair("./server.crt", "./server.key")
	if err != nil {
		log.Fatal(err)
	}

	caCert, err := ioutil.ReadFile("../ca/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{server},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAnyClientCert,
	})

	s, err := net.Listen("tcp", ":54332")
	if err != nil {
		log.Fatal(err)
	}
	gs := grpc.NewServer(grpc.Creds(ta))
	RegisterDummyServiceServer(gs, &dummyServer{})

	gs.Serve(s)
}
