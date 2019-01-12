package main

import (
	"context"
	"fmt"

	pb "github.com/uzuna/go-grpcauth/examples/pb"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }

type helloServer struct{}

func (s *helloServer) Visit(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	email := ctx.Value(CtxKeyEmail).(string)
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello %s. Echo: %s", email, req.Message),
	}, nil
}

func (s *helloServer) Authorize(ctx context.Context, fullMethodName string) error {
	// log.Println("Authorize", fullMethodName, ctx.Value(CtxKeySub).(string))
	return nil
}

func newServer() *helloServer {
	return &helloServer{}
}
