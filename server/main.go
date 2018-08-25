package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	pb "github.com/anzellai/kanosdk/kanosdk"
	uuid "github.com/satori/go.uuid"
	serial "go.bug.st/serial.v1"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Communicate(stream pb.Connector_CommunicateServer) error {
	fmt.Println("Establish connection...")
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error from %v: %s\n", request, err.Error())
			break
		}
		fmt.Printf("Stream from %s: %s\n", request.Name, request.Data)
		cReq := make(chan *pb.DeviceRequest)
		go func() {
			for r := range cReq {
				c := make(chan string)
				go connectorHandler(c, r.Name, r.Data)
				for data := range c {
					response := &pb.DeviceResponse{
						Data: string(data),
					}
					fmt.Printf("Stream sent: %s\n", data)
					stream.Send(response)
				}
			}
		}()
		cReq <- request
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

func connectorHandler(c chan string, name, data string) {
	var conn io.ReadWriteCloser
	var err error
	cache, ok := vm.Load(name)
	if ok {
		conn = cache.(io.ReadWriteCloser)
		conn.Close()
	}
	options := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
	}

	conn, err = serial.Open(name, options)
	if err != nil {
		c <- fmt.Sprintf("device error: %v", err)
		return
	}
	vm.Store(name, conn)

	bits := strings.Split(data, " ")
	u4 := uuid.Must(uuid.NewV4())
	rpcData := struct {
		Type   string   `json:"type"`
		ID     string   `json:"id"`
		Method string   `json:"method"`
		Params []string `json:"params"`
	}{
		Type:   "rpc-request",
		ID:     u4.String(),
		Method: bits[0],
		Params: []string{},
	}
	if len(bits) > 1 {
		rpcData.Params = append(rpcData.Params, bits[1:]...)
	}
	rpcByte, err := json.Marshal(rpcData)
	if err != nil {
		c <- fmt.Sprintf("marhsalling data error: %+v\n", err)
		return
	}

	rpcByte = append(rpcByte, []byte("\r\n")...)

	n, err := conn.Write(rpcByte)
	if err != nil {
		c <- fmt.Sprintf("write error: %v", err)
		return
	}
	c <- fmt.Sprintf("written %d bytes", n)

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
			result := string(response)
			response = nil
			c <- fmt.Sprintf("read %d bytes: %s\n", i, result)
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
