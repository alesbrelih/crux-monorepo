package main

import (
	"log"
	"os"
	"path"

	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/config"
	grpcAuth "github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/grpc"
	myJwt "github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/my_jwt"
	"github.com/alesbrelih/crux-monorepo/microservices/services/authentication/internal/repository"

	"google.golang.org/grpc"
)

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory. Error: %v", err)
	}
	conf := config.GetEnvConfig(cwd, "dev.env")
	start.SetUpGrpc(conf.AppPort, conf.Debug, func(server *grpc.Server) {
		// set repository and migrate db
		authRepository := repository.NewRepository(conf.DatabaseUrl)
		err = authRepository.Migrate(path.Join(cwd, "db/migrations"))
		if err != nil {
			log.Fatalf("Error migrating DB. Error: %v", err)
		}
		// set grpc controller
		jwtService := myJwt.NewJwtService(conf.JwtSecret, conf.JwtAccessExp, conf.JwtRefreshExt)
		grpcService := grpcAuth.NewAuthService(jwtService, authRepository)

		services.RegisterAuthServiceServer(server, grpcService)
	})

}
