package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/uzuna/go-grpcauth/authz"
	pb "github.com/uzuna/go-grpcauth/examples/pb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The Server Port")
)

func main() {
	// Parse flag and init options
	flag.Parse()
	listenHost := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", listenHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	authopts := grpc.UnaryInterceptor(
		// Add Authentication/Authoriza Interceptor
		grpc_middleware.ChainUnaryServer(
			grpc_auth.UnaryServerInterceptor(DebugAuthentication()),
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
