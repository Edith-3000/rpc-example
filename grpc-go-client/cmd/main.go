package main

import (
	"context"
	"log"
	"time"

	pb "grpc-go-client/proto/greet/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{
		Name: "Edith-3000",
	})
	if err != nil {
		log.Fatalf("Could not send request: %v", err)
	}

	log.Printf("Greeting from server: %s", resp.GetMessage())
}
