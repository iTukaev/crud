package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	apiValidatorPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/validator"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	botPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/update"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("Start main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := yamlPkg.MustNew()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go start(ctx, config)

	select {
	case <-c:
		log.Println("Shutting down...")
		cancel()
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
}

func start(ctx context.Context, config configPkg.Interface) {
	conn, err := grpc.Dial(config.RepoAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewUserClient(conn)

	go runGRPCServer(ctx, client, config.GRPCAddr())
	go runHTTPServer(ctx, client, config.HTTPAddr())

	runBot(ctx, client, config.BotKey())
}

func runBot(ctx context.Context, client pb.UserClient, apiKey string) {
	bot := botPkg.MustNew(apiKey)

	commandAdd := cmdAddPkg.New(client)
	bot.RegisterCommand(commandAdd)

	commandUpdate := cmdUpdatePkg.New(client)
	bot.RegisterCommand(commandUpdate)

	commandDelete := cmdDeletePkg.New(client)
	bot.RegisterCommand(commandDelete)

	commandGet := cmdGetPkg.New(client)
	bot.RegisterCommand(commandGet)

	commandList := cmdListPkg.New(client)
	bot.RegisterCommand(commandList)

	commandHelp := cmdHelpPkg.New(map[string]string{
		commandAdd.Name():    commandAdd.Description(),
		commandUpdate.Name(): commandUpdate.Description(),
		commandDelete.Name(): commandDelete.Description(),
		commandGet.Name():    commandGet.Description(),
		commandList.Name():   commandList.Description(),
	})
	bot.RegisterCommand(commandHelp)

	log.Println("Start TG bot")
	go func() {
		bot.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		bot.Stop()
		log.Println("Bot stopped")
	}
}

func runGRPCServer(ctx context.Context, client pb.UserClient, grpcSrv string) {
	listener, err := net.Listen("tcp", grpcSrv)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, apiValidatorPkg.New(client))

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

func runHTTPServer(ctx context.Context, client pb.UserClient, httpSrv string) {
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

	if err := pb.RegisterUserHandlerServer(ctx, gwMux, apiValidatorPkg.New(client)); err != nil {
		log.Fatalln("HTTP gateway register:", err)
	}

	srv := http.Server{
		Addr:    httpSrv,
		Handler: mux,
	}
	log.Println("Start HTTP gateway")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	select {
	case <-ctx.Done():
		if err := srv.Close(); err != nil {
			log.Println("HTTP server close error:", err)
		}
		log.Println("HTTP gateway stopped")
	}
}
