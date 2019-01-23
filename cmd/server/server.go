package main

import (
	"log"
	"net"
	"time"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct{}

func (s *server) Start(stream comms.Server_StartServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%v connected", req.Id)

	var i int64
	for {
		i++
		stream.Send(&comms.StartResp{Seq: i})
		_, err = stream.Recv()
		if err != nil {
			log.Println(req.Id, ":", err)
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    10 * time.Second,
		Timeout: 5 * time.Second,
	}))
	comms.RegisterServerServer(grpcServer, new(server))

	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	grpcServer.Serve(listen)
}
