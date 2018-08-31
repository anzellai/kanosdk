package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"sync"

	pb "github.com/anzellai/kanosdk/kanosdk/go"
	serial "go.bug.st/serial.v1"
	"google.golang.org/grpc"
)

type server struct{}

// Communicate is the entry point to establish device connection
func (*server) Communicate(stream pb.Connector_CommunicateServer) error {
	fmt.Println("Establish connection...")
	for {
		streamRequest, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error from %v: %s\n", streamRequest, err.Error())
			break
		}
		fmt.Printf("Stream from %s: %v\n", streamRequest.Name, streamRequest.Request)
		cReq := make(chan *pb.StreamRequest)
		go func() {
			for r := range cReq {
				c := make(chan *pb.StreamResponse)
				go connectorHandler(c, r)
				for response := range c {
					fmt.Printf("Stream sent: %v\n", response)
					stream.Send(response)
				}
			}
		}()
		cReq <- streamRequest
	}

	defer func() {
		vm.Range(func(key, value interface{}) bool {
			if cache, ok := value.(io.ReadWriteCloser); ok && cache != nil {
				fmt.Printf("Disconnecting %s...\n", key)
				cache.Close()
				return true
			}
			return false
		})
	}()

	fmt.Println("Disconnected")
	return nil
}

var vm sync.Map

func connectorHandler(c chan *pb.StreamResponse, request *pb.StreamRequest) {
	var conn io.ReadWriteCloser
	var err error
	data := &pb.StreamResponse{
		Name: request.Name,
		Response: &pb.Response{
			Type:   "rpc-response",
			Id:     request.Request.Id,
			Name:   request.Name,
			Detail: &pb.Detail{},
		},
	}
	cache, ok := vm.Load(request.Name)
	if ok {
		conn = cache.(io.ReadWriteCloser)
		conn.Close()
	}
	options := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
	}

	conn, err = serial.Open(request.Name, options)
	if err != nil {
		data.Response.Detail.Error = fmt.Sprintf("device error: %s", err.Error())
		c <- data
		return
	}
	vm.Store(request.Name, conn)

	rpcByte, err := json.Marshal(request.Request)
	if err != nil {
		data.Response.Detail.Error = fmt.Sprintf("marhsalling data error: %s", err.Error())
		c <- data
		return
	}
	rpcByte = append(rpcByte, []byte("\r\n")...)

	_, err = conn.Write(rpcByte)
	if err != nil {
		data.Response.Detail.Error = fmt.Sprintf("write error: %s", err.Error())
		c <- data
		return
	}

	var response []byte
	for {
		if conn == nil {
			return
		}
		buf := make([]byte, 1)
		i, err := conn.Read(buf)
		if err != nil {
			return
		}
		if i == 0 {
			fmt.Println("\nEOF")
			break
		}

		if string(buf) == "\n" {
			err = json.Unmarshal(response, data.Response)
			if err != nil {
				data.Response.Detail.Error = fmt.Sprintf("read error: %s", err.Error())
			} else {
				response = nil
			}
			if data.Response.Id != request.Request.Id {
				continue
			}
			c <- data
		} else {
			response = append(response, buf...)
		}
	}
}

func main() {
	port := flag.Int("port", 55555, "Server port")
	flag.Parse()

	fmt.Println("KanoSDK Connector server starting on port: ", *port)
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("Server error: %v\n", err)
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterConnectorServer(grpcServer, &server{})
	grpcServer.Serve(conn)
}
