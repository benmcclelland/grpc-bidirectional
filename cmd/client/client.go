package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8888", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := comms.NewServerClient(conn)
	stream, err := client.Start(context.Background())
	if err != nil {
		panic(err)
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	i := rand.Intn(10)
	log.Println(i)

	for j := 0; j < i; j++ {
		stream.Send(&comms.StartReq{Id: os.Args[1]})
		resp, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		log.Println("recv seq:", resp.Seq)

		time.Sleep(time.Second)
	}

	conn.Close()
}
