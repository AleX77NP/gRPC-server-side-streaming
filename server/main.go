package main

import (
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/aleksandarmilanovic/grpc-numbers/protos"
	"google.golang.org/grpc"
)

//Server ...
type server struct{
	pb.UnimplementedStreamServiceServer
}

//ConvertNumbers ...
func (s server) ConvertNumber(in *pb.Request, srv pb.StreamService_ConvertNumberServer) error {
	log.Printf("Fetch response for number: %d", in.Num);

	var wg sync.WaitGroup
	for i:=0; i<=10; i++ {
		wg.Add(1);
		go func(n int32) {
			defer wg.Done()

			time.Sleep(time.Duration(n) * time.Second)
			response := pb.Response{TextNum: convertRequestToResponse(in.Num)}
			if err := srv.Send(&response); err != nil {
				log.Printf("send error %v", err)
			}
			log.Printf("Finished request number %d", n)
		}(int32(i))
	}
	wg.Wait()
	return nil
}

func convertRequestToResponse(num int32) string {
	var res string
	if num % 2 == 0 {
		res = "Even number!"
	} else if num % 2 == 1 {
		res = "Odd number"
	} else {
		res = "Zero"
	}
	return res
}

func main() {

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen %v",err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, server{})

	log.Println("start server")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}