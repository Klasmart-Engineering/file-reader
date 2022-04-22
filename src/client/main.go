package main

import (
	"context"
	"fmt"
	"log"

	msgpb "file-reader/src/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:9000", opts...)

	if err != nil {
		log.Fatal("Failed to dial: %v", err)
	}

	defer conn.Close()

	client := msgpb.NewServiceClient(conn)

	res, err := client.Hello(context.Background(), &msgpb.Request{
		Name: "some name",
	})

	if err != nil {
		log.Fatal("Failed to GetUser: %v", err)
	}
	fmt.Printf("%+v\n", res)
}
