package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"

	"github.com/yakomisar/inventory_service/internal/docs"
	v1 "github.com/yakomisar/inventory_service/pkg/pb/inventory_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50051
	httpPort = 8081
)

type InventoryService struct {
	v1.UnimplementedInventoryServiceServer
}

func (s *InventoryService) GetItem(ctx context.Context, req *v1.GetItemRequest) (*v1.GetItemResponse, error) {
	return &v1.GetItemResponse{Name: "Item_" + req.GetId()}, nil
}

func (s *InventoryService) Ping(ctx context.Context, req *v1.PingRequest) (*v1.PingResponse, error) {
	return &v1.PingResponse{Message: "pong"}, nil
}

func main() {
	// gRPC listener
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

	// gRPC server
	grpcSrv := grpc.NewServer() // create gRPC server
	svc := &InventoryService{}  // register service
	v1.RegisterInventoryServiceServer(grpcSrv, svc)
	reflection.Register(grpcSrv) // add reflection

	// gRPC run
	go func() {
		log.Printf("gRPC server listening on %d\n", grpcPort)
		err = grpcSrv.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// HTTP gateway mux
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", grpcPort)

	if err := v1.RegisterInventoryServiceHandlerFromEndpoint(ctx, gwMux, endpoint, dialOpts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Корневой mux
	mux := http.NewServeMux()
	// REST эндпоинты REST --> gRPC
	mux.Handle("/v1/", gwMux)

	// Документация
	mux.Handle("/swagger/api.swagger.json", docs.SwaggerJSONHandler())
	mux.Handle("/docs", docs.RedocHandler())
	mux.Handle("/swagger", docs.SwaggerUIHandler())

	// Health/ready
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// CORS (разрешим локальную разработку и портал)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	httpSrv := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// HTTP run
	go func() {
		log.Printf("HTTP server (gateway+docs) listening on :%d\n", httpPort)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP serve error: %v", err)
		}
	}()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	grpcSrv.GracefulStop()
	log.Println("servers stopped gracefully")
}
