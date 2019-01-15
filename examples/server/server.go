package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/uzuna/go-grpcauth/authz"
	pb "github.com/uzuna/go-grpcauth/examples/pb"
	"google.golang.org/grpc"
)

var (
	port    = flag.Int("port", 10000, "The Server Port")
	keysURL = "https://login.microsoftonline.com/common/discovery/v2.0/keys"
)

func main() {
	// Parse flag and init options
	flag.Parse()
	listenHost := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", listenHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	atf, err := AccessTokenAuthentication(keysURL)
	_ = atf
	if err != nil {
		log.Fatalf("Fail Load JWS Keys: %v", err)
	}

	var opts []grpc.ServerOption
	authopts := grpc.UnaryInterceptor(
		// Add Authentication/Authoriza Interceptor
		grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(DebugAuthentication()),
			// grpc_auth.UnaryServerInterceptor(atf),
			authz.UnaryServerInterceptor(),
		),
	)
	opts = append(opts, authopts)
	_ = opts
	_ = authopts
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterHelloServer(grpcServer, newServer())
	log.Printf("Start listen: %s\n", listenHost)
	go func() {
		grpcServer.Serve(lis)
	}()

	code := WaitSignal()
	grpcServer.Stop()
	log.Printf("Stop Server. code: %v", code)
}
