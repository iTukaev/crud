package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	apiPkg "gitlab.ozon.dev/iTukaev/homework/internal/api"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("start gRPC server")
	user := userPkg.MustNew()

	config := yamlPkg.MustNew()
	config.Init()

	runGRPCServer(user, config.GRPCAddr())
}

func runGRPCServer(user userPkg.Interface, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, apiPkg.New(user))

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
