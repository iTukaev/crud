package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("start client")
	config := yamlPkg.MustNew()
	config.Init()

	conn, err := grpc.Dial(config.GRPCAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewUserClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	response, err := client.UserCreate(ctx, &pb.UserCreateRequest{
		Name:     "Paulo",
		Password: "123",
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("response: [%v]", response)
}
