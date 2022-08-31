package main

import (
	"context"
	"expvar"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	apiReceiverPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/receiver"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/counter"
	botPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/update"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	grpcPkg "gitlab.ozon.dev/iTukaev/homework/pkg/grpc"
	jaegerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/jaeger"
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
		_ = logger.Sync()
		cancel()
	}()

	tracer, closer, err := jaegerPkg.New(config.JService(), config.JHost())
	if err != nil {
		logger.Errorf("Jaeger initialise err: %v", err)
		return
	}
	defer func() {
		_ = closer.Close()
	}()
	opentracing.SetGlobalTracer(tracer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err = start(ctx, config, tracer, logger); err != nil {
			logger.Errorln(err)
			c <- os.Interrupt
		}
	}()

	select {
	case <-c:
	}
}

func start(ctx context.Context, config configPkg.Interface, tracer opentracing.Tracer, logger *zap.SugaredLogger) (retErr error) {
	conn, err := grpc.Dial(config.RepoAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)),
	)
	if err != nil {
		return errors.Wrap(err, "gRPC client connection")
	}

	client := pb.NewUserClient(conn)
	stopCh := make(chan struct{}, 0)

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Brokers(), cfg)
	if err != nil {
		return errors.Wrap(err, "new SyncProducer")
	}
	server := apiReceiverPkg.New(client, logger, producer)

	go func() {
		if err = runGRPCServer(ctx, server, config.GRPCAddr(), logger); err != nil {
			retErr = errors.Wrap(err, "gRPC server")
		}
		close(stopCh)
	}()
	go func() {
		if err = runHTTPServer(ctx, server, config.HTTPAddr(), logger); err != nil {
			retErr = errors.Wrap(err, "HTTP server")
		}
		close(stopCh)
	}()
	go func() {
		if err = runBot(ctx, client, config.BotKey(), logger); err != nil {
			retErr = errors.Wrap(err, "tg bot")
		}
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
	case <-stopCh:
	}
	return retErr
}

func runBot(ctx context.Context, client pb.UserClient, apiKey string, logger *zap.SugaredLogger) error {
	bot, err := botPkg.New(apiKey, logger)
	if err != nil {
		return err
	}

	commandAdd := cmdAddPkg.New(client, logger)
	bot.RegisterCommand(commandAdd)

	commandUpdate := cmdUpdatePkg.New(client, logger)
	bot.RegisterCommand(commandUpdate)

	commandDelete := cmdDeletePkg.New(client, logger)
	bot.RegisterCommand(commandDelete)

	commandGet := cmdGetPkg.New(client, logger)
	bot.RegisterCommand(commandGet)

	commandList := cmdListPkg.New(client, logger)
	bot.RegisterCommand(commandList)

	commandHelp := cmdHelpPkg.New(map[string]string{
		commandAdd.Name():    commandAdd.Description(),
		commandUpdate.Name(): commandUpdate.Description(),
		commandDelete.Name(): commandDelete.Description(),
		commandGet.Name():    commandGet.Description(),
		commandList.Name():   commandList.Description(),
	})
	bot.RegisterCommand(commandHelp)

	logger.Infoln("Start TG bot")
	stopCh := make(chan struct{}, 0)
	go func() {
		bot.Run(ctx)
		close(stopCh)
	}()

	select {
	case <-stopCh:
	case <-ctx.Done():
		bot.Stop()
	}
	logger.Infoln("Bot stopped")
	return nil
}

func runGRPCServer(ctx context.Context, server pb.UserServer, grpcSrv string, logger *zap.SugaredLogger) (retErr error) {
	listener, err := net.Listen("tcp", grpcSrv)
	if err != nil {
		return errors.Wrap(err, "listener")
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcPkg.MetricsUnaryInterceptor,
			grpcOpentracing.UnaryServerInterceptor(),
		),
		grpc.StreamInterceptor(grpcPkg.MetricsStreamInterceptor),
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

func runHTTPServer(ctx context.Context, server pb.UserServer, httpSrv string, logger *zap.SugaredLogger) (retErr error) {
	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/", gwMux)

	fs := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	mux.Handle("/counters", expvar.Handler())
	expvar.Publish("Validation service request", counter.Request)
	expvar.Publish("Validation service response", counter.Response)
	expvar.Publish("Validation service success", counter.Success)
	expvar.Publish("Validation service error", counter.Errors)

	if err := pb.RegisterUserHandlerServer(ctx, gwMux, server); err != nil {
		return errors.Wrap(err, "HTTP gateway register")
	}

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
