package main

import (
	"context"
	_ "embed"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiPkg "gitlab.ozon.dev/iTukaev/homework/internal/api"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	botPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/update"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	"gitlab.ozon.dev/iTukaev/homework/swagger"
)

func main() {
	log.Println("Start main")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := yamlPkg.MustNew()
	config.Init()

	user := userPkg.MustNew(ctx, config.PGConfig())

	go runGRPCServer(user, config.GRPCAddr())
	go runHTTPServer(config.GRPCAddr(), config.HTTPAddr())

	runBot(ctx, user, config.BotKey())
}

func runBot(ctx context.Context, user userPkg.Interface, apiKey string) {
	bot := botPkg.MustNew(apiKey)

	commandAdd := cmdAddPkg.New(user)
	bot.RegisterCommander(commandAdd)

	commandUpdate := cmdUpdatePkg.New(user)
	bot.RegisterCommander(commandUpdate)

	commandDelete := cmdDeletePkg.New(user)
	bot.RegisterCommander(commandDelete)

	commandGet := cmdGetPkg.New(user)
	bot.RegisterCommander(commandGet)

	commandList := cmdListPkg.New(user)
	bot.RegisterCommander(commandList)

	commandHelp := cmdHelpPkg.New(map[string]string{
		commandAdd.Name():    commandAdd.Description(),
		commandUpdate.Name(): commandUpdate.Description(),
		commandDelete.Name(): commandDelete.Description(),
		commandGet.Name():    commandGet.Description(),
		commandList.Name():   commandList.Description(),
	})
	bot.RegisterCommander(commandHelp)

	log.Println("Start bot")
	bot.Run(ctx)
}

func runGRPCServer(user userPkg.Interface, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, apiPkg.New(user))

	log.Println("Start gRPC")
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}

func runHTTPServer(grpcSrv, httpSrv string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterUserHandlerFromEndpoint(ctx, gwMux, grpcSrv, opts); err != nil {
		log.Fatalln(err)
	}

	mux := swagger.Mux("/swagger")
	mux.Handle("/", gwMux)

	log.Println("Start HTTP")
	if err := http.ListenAndServe(httpSrv, mux); err != nil {
		log.Fatalln(err)
	}
}
