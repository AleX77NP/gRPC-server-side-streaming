package main

import (
	"context"
	"io"
	"log"

	pb "github.com/aleksandarmilanovic/grpc-numbers/protos"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to server with error : %v",err)
	}

	client := pb.NewStreamServiceClient(conn)
	in := &pb.Request{Num: 1}
	stream, err := client.ConvertNumber(context.Background(), in)
	if err != nil {
		log.Fatalf("Cannot open stream, error: %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				done <-true
				return
			}
			if err != nil {
				log.Fatalf("Cannot recieve %v", err)
			}
			log.Printf("Response: %s", response)
		}
	}()

	<-done
	log.Printf("The end...")
}