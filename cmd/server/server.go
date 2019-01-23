package main

import (
	"log"
	"net"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"google.golang.org/grpc"
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
	grpcServer := grpc.NewServer()
	comms.RegisterServerServer(grpcServer, new(server))

	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	grpcServer.Serve(listen)
}
