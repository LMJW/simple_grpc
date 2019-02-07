package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cred, err := credentials.NewClientTLSFromFile("./ca.crt", "")
	if err != nil {
		log.Fatal(err)
	}

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
