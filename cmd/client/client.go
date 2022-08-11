package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
)

func main() {
	log.Println("start client")
	config := yamlPkg.MustNew()

	conn, err := grpc.Dial(config.GRPCAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewUserClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	response, err := client.UserAllList(ctx, &pb.UserAllListRequest{
		Order: true,
		Limit: 5,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		next, err := response.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		for i, user := range next.Users {
			fmt.Println(i, user.String())
		}
	}

}
