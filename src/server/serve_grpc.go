package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

func ServeGrpc(ctx context.Context, addr string) (err error) {
	defer log.Println("GRPC server stopped", err)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	server := grpc.NewServer()

	log.Printf("Listen and Serve Request-Grpc-Service API at: %s\n", addr)
	return server.Serve(lis)
}
