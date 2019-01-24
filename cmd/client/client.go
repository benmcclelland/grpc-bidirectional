package main

import (
	"log"
	"os"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial err :%v", err)
	}
	defer conn.Close()

	client := comms.NewWorkClient(conn)
	stream, err := client.Hello(context.Background())
	if err != nil {
		log.Fatalf("start err :%v", err)
	}

	stream.Send(&comms.Req{Id: os.Args[1]})

	c := make(chan *comms.Resp)
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Fatalf("recv err :%v", err)
			}
			c <- resp
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			log.Fatalf("server disconnected")
		case r := <-c:
			log.Println("recv work:", r.Seq)
		}
	}
}
