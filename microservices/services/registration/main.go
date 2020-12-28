package main

import (
	"log"
	"os"

	envConf "github.com/alesbrelih/crux-monorepo/microservices/config"
	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/internal/repository"

	"google.golang.org/grpc"
)

func main() {

	// TODO: this isnt right -> move config to diff place -> simplyfy this
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory. Error: %v", err)
	}
	log.Println(cwd)
	envConf := envConf.GetEnvConfig(cwd, "dev.env")

	start.SetUpGrpc(envConf, func(server *grpc.Server) {

		repo := repository.NewRepository(envConf.DatabaseUrl)
		grpcService := grpcAuth.NewAuthService(jwtService, authRepository)

		services.RegisterRegistrationServiceServer(server, grpcService)
	})

}
