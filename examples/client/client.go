package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/joho/godotenv"
	"github.com/uzuna/go-grpcauth/examples/oauthclient"
	pb "github.com/uzuna/go-grpcauth/examples/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var serverAddr = flag.String("server_addr", "localhost:10000", "The server address in format of [host:port]")

func main() {
	flag.Parse()

	// Start Collect Token
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := oauthclient.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var token string
	if len(config.Token) < 1 {
		token, err = oauthclient.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		token = config.Token
	}

	// start GRPC
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
	log.Println(token)
	ctx = ctxWithToken(ctx, "bearer", token)

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
