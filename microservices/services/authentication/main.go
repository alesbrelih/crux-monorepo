package main

import (
	"log"
	"os"
	"path"

	envConf "github.com/alesbrelih/crux-monorepo/microservices/config"
	"github.com/alesbrelih/crux-monorepo/microservices/internal/start"
	"github.com/alesbrelih/crux-monorepo/microservices/protos/build/services"
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
	log.Println(cwd)
	envConf := envConf.GetEnvConfig(cwd, "dev.env")

	start.SetUpGrpc(envConf, func(server *grpc.Server) {
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
	})

}
