package main

import (
	"fmt"
	golog "log"
	"net"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/caarlos0/env/v6"
	grpcServiceImpl "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/endpoints"
	proto "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/pb"
	service "github.com/cfloress-gb-cl/final-project-bootcamp/cmd/grpc/user/service"
	mongo "github.com/cfloress-gb-cl/final-project-bootcamp/repository/mongodb"
	mysql "github.com/cfloress-gb-cl/final-project-bootcamp/repository/mysql"
	domain "github.com/cfloress-gb-cl/final-project-bootcamp/repository/user"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	// if we crash the go code, we get the file name and line number
	golog.SetFlags(golog.LstdFlags | golog.Lshortfile)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))

	if err != nil {
		panic(fmt.Sprintf("Could not create the listener %v", err))
	}

	userService := domain.NewUserService(getActiveRepository())
	endpoints := grpcServiceImpl.NewGrpcUsersServer(userService)

	grpcUserServer := service.NewGrpcUserServer(*endpoints, logger)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	proto.RegisterUsersServer(baseServer, grpcUserServer)
	fmt.Println("grpc server started!..")

	if err := baseServer.Serve(ls); err != nil {
		panic(fmt.Sprintf("failed to serve: %s", err))
	}
}

func getActiveRepository() domain.Repository {

	envVar := os.Getenv("USERS_REPOSITORY")

	fmt.Println(envVar)

	if len(envVar) == 0 {
		envVar = "mongo"
	}

	switch envVar {
	case "mongo":
		repo, err := mongo.NewMongoUserRepository()
		if err != nil {
			panic(fmt.Sprintf("mongoDB connection failed: %s", err))
		}
		return repo
	case "mysql":
		repo, err := mysql.NewMySQLUserRepository()
		if err != nil {
			panic(fmt.Sprintf("mysql connection failed: %s", err))
		}
		return repo
	}
	return nil
}

type config struct {
	Port int `env:"GRPCSERVICE_PORT" envDefault:"9000"`
}
