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
	config, _ := yamlPkg.New()

	conn, err := grpc.Dial(config.GRPCAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewUserClient(conn)
	//{
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "123456789")
	//
	//	res, err := client.UserGet(ctx, &pb.UserGetRequest{Name: "Piter"})
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(res)
	//}
	//time.Sleep(1 * time.Second)
	//{
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "987654321")
	//	resCr, err := client.UserCreate(ctx, &pb.UserCreateRequest{
	//		User: &pbModels.User{
	//			Name:     "Timut",
	//			Password: "pass",
	//			Email:    "tm@il.ru",
	//			FullName: "Tim Owner",
	//		},
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(resCr)
	//}
	//time.Sleep(1 * time.Second)
	//{
	//	pass, email, full := "pass111", "tm@il.ru11", "Tim Owner"
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "000111000")
	//	resCr, err := client.UserUpdate(ctx, &pb.UserUpdateRequest{
	//		Name: "Timut",
	//		Profile: &pbModels.Profile{
	//			Password: &pass,
	//			Email:    &email,
	//			FullName: &full,
	//		},
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(resCr)
	//}
	//time.Sleep(1 * time.Second)
	{
		ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "222333222")
		resCr, err := client.UserDelete(ctx, &pb.UserDeleteRequest{
			Name: "Timut",
		})
		if err != nil {
			log.Println(err)
		}
		log.Println(resCr)
	}
	//time.Sleep(1 * time.Second)
	//{
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "222333222")
	//	resCr, err := client.UserList(ctx, &pb.UserListRequest{
	//		Order:  true,
	//		Limit:  3,
	//		Offset: 0,
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(resCr)
	//}
	//time.Sleep(1 * time.Second)
	//{
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "444555666")
	//	response, err := client.UserAllList(ctx, &pb.UserAllListRequest{
	//		Order: false,
	//		Limit: 2,
	//	})
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	for {
	//		next, err := response.Recv()
	//		if errors.Is(err, io.EOF) {
	//			break
	//		}
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		for i, user := range next.Users {
	//			fmt.Println(i, user.String())
	//		}
	//	}
	//}
}
