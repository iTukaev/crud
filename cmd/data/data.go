package main

import (
	"context"
	"expvar"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	apiDataPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/data"
	dataPkg "gitlab.ozon.dev/iTukaev/homework/internal/brokers/data"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	repoPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo"
	localCachePkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/local"
	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	jaegerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/jaeger"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
	redisPkg "gitlab.ozon.dev/iTukaev/homework/pkg/redis"
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
			workers = 10
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

	client, err := redisPkg.New(ctx, config.RedisConfig())
	if err != nil {
		return errors.Wrap(err, "new redis client")
	}

	user := userPkg.New(data, logger, client)

	tracer, closer, err := jaegerPkg.New(config.JService(), config.JHost())
	if err != nil {
		logger.Errorf("Jaeger initialise err: %v", err)
		return
	}
	defer func() {
		_ = closer.Close()
	}()
	opentracing.SetGlobalTracer(tracer)

	server := apiDataPkg.New(user, logger)

	stopCh := make(chan struct{}, 0)
	go func() {
		if err = runGRPCServer(ctx, server, config.GRPCDataAddr(), logger); err != nil {
			retErr = errors.Wrap(err, "gRPC server")
		}
		close(stopCh)
	}()
	go func() {
		if err = runHTTPServer(ctx, config.HTTPDataAddr(), logger); err != nil {
			retErr = errors.Wrap(err, "HTTP server")
		}
		close(stopCh)
	}()
	go func() {
		if err = runService(ctx, config.Brokers(), logger, user); err != nil {
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

func runGRPCServer(ctx context.Context, server pb.UserServer, grpcSrv string, logger *zap.SugaredLogger) (retErr error) {
	listener, err := net.Listen("tcp", grpcSrv)
	if err != nil {
		log.Fatalln("Listener create:", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcOpentracing.UnaryServerInterceptor()),
		grpc.StreamInterceptor(grpcOpentracing.StreamServerInterceptor()),
	)
	pb.RegisterUserServer(grpcServer, server)

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

	handler := dataPkg.NewHandler(user, logger, producer)

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

func runHTTPServer(ctx context.Context, httpSrv string, logger *zap.SugaredLogger) (retErr error) {
	mux := http.NewServeMux()
	mux.Handle("/counters", expvar.Handler())
	expvar.Publish("Hit cache", counter.Hit)
	expvar.Publish("Miss cache", counter.Miss)

	srv := http.Server{
		Addr:    httpSrv,
		Handler: mux,
	}
	logger.Infoln("Start HTTP gateway")
	stopCh := make(chan struct{}, 0)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			retErr = errors.Wrap(err, "ListenAndServe")
		}
		close(stopCh)
	}()

	defer func() {

	}()
	select {
	case <-stopCh:
	case <-ctx.Done():
		if err := srv.Close(); err != nil {
			logger.Errorln("HTTP server close error:", err)
		}
	}
	logger.Infoln("HTTP gateway stopped")
	return nil
}
