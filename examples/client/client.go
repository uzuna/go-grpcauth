package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	pb "github.com/uzuna/go-grpcauth/examples/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var serverAddr = flag.String("server_addr", "localhost:10000", "The server address in format of [host:port]")

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// grpc metadataにAuthorization: Bearer <Token>を与える
	ctx = ctxWithToken(ctx, "bearer", "test")

	req := &pb.HelloRequest{
		Message: "Mime",
	}
	rep, err := client.Visit(ctx, req)
	log.Println(rep, err)
}

// Authorization にTokenを入れる
func ctxWithToken(ctx context.Context, scheme string, token string) context.Context {
	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", scheme, token))
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)
	return nCtx
}
