package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
)

func main() {
	log.Println("start client")
	config, _ := yamlPkg.New()

	conn, err := grpc.Dial(config.GRPCAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewUserClient(conn)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	_, err = client.UserCreate(ctx, &pb.UserCreateRequest{
		User: &pbModels.User{
			Name:     "IA",
			Password: "123",
			Email:    "123@123.ru",
			FullName: "oslik",
		},
	})
	if err != nil {
		log.Println(err)
	}
	//response, err := client.UserAllList(ctx, &pb.UserAllListRequest{
	//	Order: false,
	//	Limit: 2,
	//})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//for {
	//	next, err := response.Recv()
	//	if errors.Is(err, io.EOF) {
	//		break
	//	}
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	for i, user := range next.Users {
	//		fmt.Println(i, user.String())
	//	}
	//}
}
