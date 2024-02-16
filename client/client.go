package main

import (
	"context"
	"fmt"
	"main/grpcerror"
	pb "main/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	client := pb.NewErrorTestClient(conn)

	_, err = client.ErrorTest(context.Background(), &pb.EmptyMsg{})
	if err != nil {
		protoErr := grpcerror.NewErorFromProto(err)
		fmt.Println(protoErr.Error())
	}
}
