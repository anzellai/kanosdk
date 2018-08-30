package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/satori/go.uuid"

	pb "github.com/anzellai/kanosdk/kanosdk"
	"google.golang.org/grpc"
)

func main() {
	address := flag.String("address", "localhost:55555", "Server address to connect to")
	flag.Parse()
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewConnectorClient(conn)
	stream, err := client.Communicate(context.Background())
	if err != nil {
		log.Printf("Connection error: %v", err)
		return
	}
	waitc := make(chan struct{})
	devices := map[string]string{}
	go func() {
		device, err := devicePrompt()
		if err != nil {
			return
		}
		log.Println("Communication is bidirectional, send WRITE more than once, READ stream will return concurrently.")
		devices[device] = uuid.NewV4().String()
		for {
			data, err := commPrompt()
			if err != nil {
				break
			}
			if data == "q" || data == "quit" {
				log.Println("QUITTING...")
				waitc <- struct{}{}
			}
			bits := strings.Split(data, ":")
			req := &pb.StreamRequest{
				Name: device,
				Request: &pb.Request{
					Type:   "rpc-request",
					Id:     devices[device],
					Method: bits[0],
				},
			}
			if len(bits) > 1 {
				req.Request.Params = []*pb.Param{&pb.Param{Mode: bits[1]}}
			}

			log.Printf("Sending: %v\n", req)
			stream.Send(req)
		}
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				log.Printf("Receiving error: %v", err)
				return
			}
			log.Printf("Received: %+v", response)
		}
	}()

	<-waitc
	stream.CloseSend()
}
