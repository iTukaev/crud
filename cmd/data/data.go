package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	apiDataPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/data"
	dataPkg "gitlab.ozon.dev/iTukaev/homework/internal/brokers/data"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
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
		if err = start(ctx, config, logger); err != nil {
			logger.Errorln("gRPC", err)
		}
		c <- os.Interrupt
	}()

	<-c
}

func start(ctx context.Context, config configPkg.Interface, logger *zap.SugaredLogger) (retErr error) {
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

	stopCh := make(chan struct{}, 0)
	go func() {
		if err := runGRPCServer(ctx, config.GRPCAddr(), logger, user); err != nil {
			retErr = errors.Wrap(err, "gRPC server")
		}
		close(stopCh)
	}()
	go func() {
		if err := runService(ctx, config.Brokers(), logger, user); err != nil {
			retErr = errors.Wrap(err, "consumer service")
		}
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
	case <-stopCh:
	}
	return retErr
}

func runGRPCServer(ctx context.Context, grpcSrv string, logger *zap.SugaredLogger, user userPkg.Interface) (retErr error) {
	listener, err := net.Listen("tcp", grpcSrv)
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

func runService(ctx context.Context, brokers []string, logger *zap.SugaredLogger, user userPkg.Interface) error {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return errors.Wrap(err, "new SyncProducer")
	}

	income, err := sarama.NewConsumerGroup(brokers, consts.GroupData, cfg)
	if err != nil {
		return errors.Wrap(err, "new ConsumerGroup")
	}

	handler := dataPkg.NewHandler(ctx, user, logger, producer)

	go func() {
		for {
			if err = income.Consume(ctx, []string{consts.TopicData}, handler); err != nil {
				logger.Errorf("on consume: <%v>", err)
				time.Sleep(time.Second * 5)
			}
		}
	}()

	<-ctx.Done()
	return income.Close()
}
