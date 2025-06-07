package main

import (
	"context"
	"fmt"
	"io"
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

	doUnary(client)
	doServerStreaming(client)
	doClientStreaming(client)
	doBiDiStreaming(client)
}

// 1. Unary RPC
func doUnary(client pb.GreeterClient) {
	fmt.Println("🔹 Unary RPC:")
	req := &pb.HelloRequest{Name: "Kapil"}
	res, err := client.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalf("Error in SayHello: %v", err)
	}
	fmt.Println("Response:", res.Message)
}

// 2. Server Streaming RPC
func doServerStreaming(client pb.GreeterClient) {
	fmt.Println("\n🔹 Server Streaming RPC:")
	req := &pb.HelloRequest{Name: "Kapil"}
	stream, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error in GreetManyTimes: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream recv error: %v", err)
		}
		fmt.Println("Received:", msg.Message)
	}
}

// 3. Client Streaming RPC
func doClientStreaming(client pb.GreeterClient) {
	fmt.Println("\n🔹 Client Streaming RPC:")
	stream, err := client.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error starting LongGreet: %v", err)
	}

	names := []string{"Alice", "Bob", "Charlie"}
	for _, name := range names {
		fmt.Println("Sending:", name)
		stream.Send(&pb.HelloRequest{Name: name})
		time.Sleep(500 * time.Millisecond)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving LongGreet response: %v", err)
	}
	fmt.Println("Response:", reply.Message)
}

// 4. Bidirectional Streaming RPC
func doBiDiStreaming(client pb.GreeterClient) {
	fmt.Println("\n🔹 Bidirectional Streaming RPC:")
	stream, err := client.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error starting GreetEveryone: %v", err)
	}

	waitChan := make(chan struct{})

	// Receive
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Receive error: %v", err)
			}
			fmt.Println("Received:", resp.Message)
		}
		close(waitChan)
	}()

	// Send
	names := []string{"Amar", "Akbar", "Anthony"}
	for _, name := range names {
		fmt.Println("Sending:", name)
		stream.Send(&pb.HelloRequest{Name: name})
		time.Sleep(500 * time.Millisecond)
	}

	stream.CloseSend()
	<-waitChan
}
