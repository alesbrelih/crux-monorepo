package main

import (
	"log"
	"os"

	envConf "github.com/alesbrelih/crux-monorepo/microservices/config"
	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
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
	log.Println(cwd)
	envConf := envConf.GetEnvConfig(cwd, "dev.env")

	start.SetUpGrpc(envConf, func(server *grpc.Server) {

		repo := repository.NewRepository(envConf.DatabaseUrl)
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "registration",
			Level: hclog.LevelFromString(envConf.LogLevel), // todo set from env
		})
		grpcService := grpcRegistration.NewRegistrationService(logger, repo)

		services.RegisterRegistrationServiceServer(server, grpcService)
	})

}
