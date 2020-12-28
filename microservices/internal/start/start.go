package start

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/alesbrelih/crux-monorepo/microservices/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpGrpc(envConf *config.Enviroment, setUp func(*grpc.Server)) {

	// default listener
	lis, err := net.Listen("tcp", envConf.AppPort)
	if err != nil {
		log.Fatalf("Failed to listen: '%v' on port '%v'", err, envConf.AppPort)
	}

	// set grpc server
	server := grpc.NewServer()

	// enable reflection on debug
	// this allows service discovery
	if envConf.Debug {
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

	// set up specific routes for current grpc service
	setUp(server)

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
