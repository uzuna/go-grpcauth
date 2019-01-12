package main

import (
	"context"
	"fmt"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/pkg/errors"
	pb "github.com/uzuna/go-grpcauth/examples/pb"
)

type helloServer struct{}

func (s *helloServer) Visit(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.Message),
	}, nil
}

func (s *helloServer) Authorize(ctx context.Context, fullMethodName string) error {
	log.Println("Authorize", fullMethodName)
	return nil
}

func newServer() *helloServer {
	return &helloServer{}
}

// DebugAuthentication is 認証情報確認をする
func DebugAuthentication() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		log.Println("Call Auth Method")
		// Contextから取り出した値を検証する
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		log.Println("Auth", token)
		return ctx, nil
	}
}
