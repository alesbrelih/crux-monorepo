package main

import (
	"log"
	"os"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/pkg"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/registration/config"
	grpcUser "github.com/alesbrelih/crux-monorepo/microservices/services/user/internal/grpc"
	"github.com/alesbrelih/crux-monorepo/microservices/services/user/internal/repository"
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
		passUtil := pkg.NewPasswordUtil()
		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "user",
			Level: hclog.LevelFromString(config.LogLevel), // todo set from env
		})
		grpcService := grpcUser.NewUserService(logger, passUtil, repo)

		services.RegisterUserServiceServer(server, grpcService)
	})

}
