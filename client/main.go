package main

import (
	"context"
	"flag"
	"log"

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
	waitc := make(chan struct{})
	go func() {
		device, err := devicePrompt()
		if err != nil {
			return
		}
		log.Println("Communication is bidirectional, send WRITE more than once, READ stream will return concurrently.")
		for {
			data, err := commPrompt()
			if err != nil {
				break
			}
			if data == "q" || data == "quit" {
				log.Println("QUITTING...")
				waitc <- struct{}{}
			}
			req := &pb.DeviceRequest{
				Name: device,
				Data: data,
			}

			log.Printf("Sending: %v\n", req)
			stream.Send(req)
		}
	}()

	<-waitc
	stream.CloseSend()
}
