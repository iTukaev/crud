package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	apiPkg "gitlab.ozon.dev/iTukaev/homework/internal/api"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	botPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/update"
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

	start(ctx, config)
}

func start(ctx context.Context, config configPkg.Interface) {
	var data repoPkg.Interface
	if config.Local() {
		data = localCachePkg.New(config.WorkersCount())
	} else {
		pg := config.PGConfig()
		data = postgresPkg.MustNew(ctx, pg.Host, pg.Port, pg.User, pg.Password, pg.DBName)
	}
	user := userPkg.MustNew(data)

	go runGRPCServer(user, config.GRPCAddr())
	go runHTTPServer(user, config.HTTPAddr())

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

	log.Printf("Start TG bot")
	bot.Run(ctx)
}

func runGRPCServer(user userPkg.Interface, grpcSrv string) {
	listener, err := net.Listen("tcp", grpcSrv)
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

func runHTTPServer(user userPkg.Interface, httpSrv string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	if err := pb.RegisterUserHandlerServer(ctx, gwMux, apiPkg.New(user)); err != nil {
		log.Fatalln("HTTP gateway register:", err)
	}

	log.Println("Start HTTP gateway")
	if err := http.ListenAndServe(httpSrv, mux); err != nil {
		log.Fatalln(err)
	}
}
