package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"main/grpcerror"
	pb "main/proto"
)

type ErrorTestServer struct {
	pb.UnimplementedErrorTestServer
}

func (s ErrorTestServer) ErrorTest(ctx context.Context, _ *pb.EmptyMsg) (*pb.StatusMsg, error) {
	err := grpcerror.NewExternal().
		WithDescription("WOW! what a day").
		WithCode("test error")

	return &pb.StatusMsg{}, err
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterErrorTestServer(grpcServer, ErrorTestServer{})

	if err := grpcServer.Serve(lis); err != nil {
		fmt.Print("failed to start server")
	}
}
