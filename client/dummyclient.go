package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func main() {
	cp := x509.NewCertPool()
	bs, err := ioutil.ReadFile("./ca.crt")
	cp.AppendCertsFromPEM(bs)

	clientcert, err := tls.LoadX509KeyPair("./client.crt", "./client.key")

	cfg := tls.Config{
		Certificates: []tls.Certificate{clientcert},
		RootCAs:      cp,
	}

	cred := credentials.NewTLS(&cfg)

	con, err := grpc.Dial("localhost:54332", grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatal(err)
	}

	cc := NewDummyServiceClient(con)

	ch := make(chan os.Signal)
	defer close(ch)

	go func() {
		i := 0
		for {
			r, err := cc.GetDummy(context.Background(), &DummyRequest{Count: int32(i), Msg: "send"})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(r)
			time.Sleep(5 * time.Second)
		}
	}()
	select {
	case <-ch:
		return
	}
}
