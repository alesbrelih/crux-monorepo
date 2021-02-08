package main

import (
	"log"
	"os"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	grpcService "github.com/alesbrelih/crux-monorepo/microservices/services/mail/internal/grpc"
	"github.com/alesbrelih/crux-monorepo/microservices/services/mail/internal/repository"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/config"
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
			Name:  "mail",
			Level: hclog.LevelFromString(config.LogLevel), // todo set from env
		})

		grpcService := grpcService.NewMailService(logger, repo)

		services.RegisterMailServiceServer(server, grpcService)
	})

}
