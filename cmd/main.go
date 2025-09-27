package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	v1 "github.com/yakomisar/inventory_service/pkg/pb/inventory_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

type InventoryService struct {
	v1.UnimplementedInventoryServiceServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// create gRPC server
	srv := grpc.NewServer()

	// register service
	svc := &InventoryService{}

	v1.RegisterInventoryServiceServer(srv, svc)

	// add reflection
	reflection.Register(srv)

	// launch gRPC server
	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = srv.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// launch http server that includes gRPC Gateway and Swagger UI
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// create mux for http requests
		mux := runtime.NewServeMux()
	}()
}
