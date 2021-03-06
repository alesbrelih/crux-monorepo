package main

import (
	"log"
	"os"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/grpc_clients"
	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/config"
	grpcRegistration "github.com/alesbrelih/crux-monorepo/microservices/services/registration/internal/grpc"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/internal/repository"
	"github.com/hashicorp/go-hclog"

	"google.golang.org/grpc"
)

func main() {

	// TODO: this isnt right -> move config to diff place -> simplyfy this
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory. Error: %v", err)
	}
	config := config.GetEnvConfig(cwd, "dev.env")

	start.SetUpGrpc(config.AppPort, config.Debug, func(server *grpc.Server) {

		repo := repository.NewRepository(config.DatabaseUrl)
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "registration",
			Level: hclog.LevelFromString(config.LogLevel), // todo set from env
		})
		userClient := grpc_clients.NewUserClient("localhost:9010")
		grpcService := grpcRegistration.NewRegistrationService(logger, repo, userClient)

		services.RegisterRegistrationServiceServer(server, grpcService)
	})

}
