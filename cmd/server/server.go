package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/benmcclelland/grpc-bidirectional/comms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	sync.Mutex
	clients map[string]comms.Work_HelloServer
}

func (s *server) Hello(stream comms.Work_HelloServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Println(err)
		return err
	}

	s.Lock()
	s.clients[req.Id] = stream
	s.Unlock()

	log.Printf("%v connected", req.Id)

	select {
	case <-stream.Context().Done():
		log.Printf("%v disconnected", req.Id)
		s.Lock()
		delete(s.clients, req.Id)
		s.Unlock()
		return nil
	}
}

func main() {
	var port string
	flag.StringVar(&port, "p", ":8888", "server")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    10 * time.Second,
		Timeout: 5 * time.Second,
	}))

	s := &server{
		clients: make(map[string]comms.Work_HelloServer),
	}

	comms.RegisterWorkServer(grpcServer, s)

	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	go grpcServer.Serve(listen)

	var i int64
	for {
		s.Lock()
		if len(s.clients) == 0 {
			s.Unlock()
			time.Sleep(time.Second)
			continue
		}

		i++
		fmt.Println("connected clients:")
		for k, v := range s.clients {
			fmt.Println(k)
			err := v.Send(&comms.Resp{Seq: i})
			if err != nil {
				log.Printf("send %v: %v", k, err)
			}
		}
		s.Unlock()
		time.Sleep(time.Second)
	}
}
