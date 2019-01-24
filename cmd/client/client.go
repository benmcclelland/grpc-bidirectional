package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	var server string
	flag.StringVar(&server, "s", "127.0.0.1:8888", "server")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if len(flag.Args()) != 1 {
		log.Fatalf("must supply client name: ./client <name>")
	}

	fmt.Println("connecting to:", server)

	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial err :%v", err)
	}
	defer conn.Close()

	client := comms.NewWorkClient(conn)
	var stream comms.Work_HelloClient
	for {
		stream, err = client.Hello(context.Background())
		if err == nil {
			break
		}
		if err != nil {
			errStatus, _ := status.FromError(err)
			if codes.Unavailable != errStatus.Code() {
				log.Fatalf("start err :%v", err)
			}
		}
	}

	err = stream.Send(&comms.Req{Id: flag.Args()[0]})
	if err != nil {
		log.Fatalf("send err: %v", err)
	}
	err = stream.CloseSend()
	if err != nil {
		log.Fatalf("closesend err: %v", err)
	}

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
