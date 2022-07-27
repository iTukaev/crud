package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("start http server")

	config := yamlPkg.MustNew()
	config.Init()

	runHTTPServer(config)
}

func runHTTPServer(config configPkg.Interface) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterUserHandlerFromEndpoint(ctx, mux, config.GRPCAddr(), opts); err != nil {
		log.Fatalln(err)
	}

	if err := http.ListenAndServe(config.HTTPAddr(), mux); err != nil {
		log.Fatalln(err)
	}
}
