package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port  = ":9000"
	debug = true
)

func main() {

	// default listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: '%v' on port '%v'", err, port)
	}

	// set grpc server
	server := grpc.NewServer()

	// enable reflection on debug
	// this allows service discovery
	if debug {
		reflection.Register(server)
	}

	// server error
	errChan := make(chan error)

	// stop sign
	stopChan := make(chan os.Signal)

	// stop signals
	// SIGINT -> ctrl+C
	// SIGTERM default signal sent to kill
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	services.RegisterAuthServiceServer(server, services.UnimplementedAuthServiceServer{})

	// start server in separate goroutine so its nonblocking
	go func() {
		if err := server.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	// defer graceful shutdown
	defer func() {
		server.GracefulStop()
	}()

	// wait/block for any sigterm signal
	// once recieved defered graceful shutdown will start
	select {
	case err := <-errChan:
		log.Printf("Fatal error starting grpc server: %v", err)
	case <-stopChan:
		log.Printf("Recieved shutdown signal. Starting graceful shutdown.")
	}

}
