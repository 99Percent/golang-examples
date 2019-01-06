// +build js,wasm

package main

import (
	"context"
	"fmt"
	"log"
	"syscall/js"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/mycodesmells/golang-examples/grpc/proto/service"
)

func main() {
	fmt.Println("Hello World!")

	registerCallbacks()
	c := make(chan struct{})
	<-c
}

func registerCallbacks() {
	js.Global().Set("callGRPC", js.NewCallback(callGRPC))
}

func callGRPC(i []js.Value) {
	conn, err := grpc.Dial("http://localhost:6000", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to start gRPC connection: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewSimpleServerClient(conn)

	md := metadata.Pairs("token", "valid-token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = client.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return
	}
	log.Println("Created user!")

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Username: "slomek"})
	if err != nil {
		log.Printf("Failed to get created user: %v", err)
		return
	}
	log.Printf("User exists: %v\n", resp)

	resp2, err := client.GreetUser(ctx, &pb.GreetUserRequest{Greeting: "howdy", Username: "slomek"})
	if err != nil {
		log.Printf("Failed to greet user: %v", err)
		return
	}
	log.Printf("Greeting: %s\n", resp2.Greeting)
}
