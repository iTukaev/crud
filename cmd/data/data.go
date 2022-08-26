package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	apiDataPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/data"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	localCachePkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/local"
	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

func main() {
	config, err := yamlPkg.New()
	if err != nil {
		log.Fatalln("Config init error:", err)
	}
	logger, err := loggerPkg.New(config.LogLevel())
	if err != nil {
		log.Fatalln("Config init error:", err)
	}
	logger.Infoln("Start main")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Infoln("Shutting down...")
		cancel()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err = runGRPCServer(ctx, config, logger); err != nil {
			logger.Errorln("gRPC", err)
		}
		c <- os.Interrupt
	}()

	<-c
}

func runGRPCServer(ctx context.Context, config configPkg.Interface, logger *zap.SugaredLogger) (retErr error) {
	var data repoPkg.Interface
	if config.Local() {
		workers := config.WorkersCount()
		if workers == 0 {
			workers = runtime.NumCPU()
		}
		data = localCachePkg.New(workers, logger)
	} else {
		pg := config.PGConfig()
		pool, err := postgresPkg.NewPostgres(ctx, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName, logger)
		if err != nil {
			logger.Errorln("New Postgres", err)
			return err
		}
		data = postgresPkg.New(pool, logger)
	}
	user := userPkg.New(data, logger)

	listener, err := net.Listen("tcp", config.RepoAddr())
	if err != nil {
		log.Fatalln("Listener create:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, apiDataPkg.New(user, logger))

	logger.Infoln("Start gRPC")

	stopCh := make(chan struct{}, 0)
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			retErr = errors.Wrap(err, "serve")
		}
		close(stopCh)
	}()

	select {
	case <-stopCh:
	case <-ctx.Done():
		grpcServer.Stop()
	}
	logger.Infoln("gRPC stopped")
	return
}
