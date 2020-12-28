package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	envConf "github.com/alesbrelih/crux-monorepo/microservices/services/authentication/config"
	grpcAuth "github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/grpc"
	myJwt "github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/my_jwt"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/repository"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port  = ":9000"
	debug = true
)

func getEnvConfig(path, file string) *envConf.Enviroment {
	// get enviroment variables
	var envConf envConf.Enviroment
	viper.SetConfigFile(file)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read enviroment config. Err: %v", err)
	}
	if err := viper.Unmarshal(&envConf); err != nil {
		log.Fatalf("Failed to unmarshal enviroment variables. Err: %v", err)
	}
	return &envConf
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory. Error: %v", err)
	}
	log.Println(cwd)
	envConf := getEnvConfig(cwd, "dev.env")

	// default listener
	lis, err := net.Listen("tcp", envConf.AppPort)
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

	// set repository and migrate db
	authRepository := repository.NewRepository(envConf.DatabaseUrl)
	err = authRepository.Migrate(path.Join(cwd, "db/migrations"))
	if err != nil {
		log.Fatalf("Error migrating DB. Error: %v", err)
	}
	// set grpc controller

	jwtService := myJwt.NewJwtService(envConf.JwtSecret, envConf.JwtAccessExp, envConf.JwtRefreshExt)
	grpcService := grpcAuth.NewAuthService(jwtService, authRepository)

	services.RegisterAuthServiceServer(server, grpcService)

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
