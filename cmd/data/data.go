package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	apiDataPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/data"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	localCachePkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/local"
	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("Start main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := yamlPkg.MustNew()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go runGRPCServer(ctx, config)

	select {
	case <-c:
		log.Println("Shutting down...")
		cancel()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
}

func runGRPCServer(ctx context.Context, config configPkg.Interface) {
	var data repoPkg.Interface
	if config.Local() {
		data = localCachePkg.New(config.WorkersCount())
	} else {
		pg := config.PGConfig()
		data = postgresPkg.MustNew(ctx, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName)
	}
	user := userPkg.MustNew(data)

	listener, err := net.Listen("tcp", config.RepoAddr())
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, apiDataPkg.New(user))

	log.Println("Start gRPC")

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalln(err)
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.Stop()
		log.Println("gRPC stopped")
	}
}
