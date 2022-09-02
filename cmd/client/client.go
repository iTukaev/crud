package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	redisPkg "gitlab.ozon.dev/iTukaev/homework/pkg/redis"
)

func main() {
	log.Println("start client")
	config, _ := yamlPkg.New()

	ctx := context.Background()
	conn, err := grpc.Dial(config.GRPCAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	redisCl, err := redisPkg.New(ctx, config.RedisConfig())
	if err != nil {
		log.Println("redis", err)
		return
	}
	client := pb.NewUserClient(conn)
	{
		wg := sync.WaitGroup{}
		wg.Add(1)
		pub := pb.Wait_pub
		if pub == pb.Wait_pub {
			go func() {
				defer wg.Done()
				pubSub := redisCl.Subscribe(ctx, "get")
				defer pubSub.Close()
				msg, err := pubSub.ReceiveMessage(ctx)
				if err != nil {
					log.Println("ReceiveMessage", err)
				}
				var user models.User
				if err = json.Unmarshal([]byte(msg.Payload), &user); err != nil {
					log.Println("unmarshal", err)
				}
				fmt.Println(user)

			}()
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "meta", "123456789")

		res, err := client.UserGet(ctx, &pb.UserGetRequest{
			Name:   "Piter",
			PubSub: pub,
		})
		if err != nil {
			log.Println(err)
		}
		log.Println(res)

		switch pub {
		case pb.Wait_cache:
			go func() {
				defer wg.Done()
				time.Sleep(1 * time.Second)
				resp, err := client.Data(ctx, &pb.DataRequest{Uid: res.GetUid()})
				if err != nil {
					log.Println("data", err)
					return
				}

				var user models.User
				if err = json.Unmarshal(resp.Body.Value, &user); err != nil {
					log.Println("unmarshal", err)
				}
				fmt.Println(user)
			}()
		}
		wg.Wait()
	}

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
	//{
	//	ctx := metadata.AppendToOutgoingContext(context.Background(), "meta", "222333222")
	//	resCr, err := client.UserDelete(ctx, &pb.UserDeleteRequest{
	//		Name: "Timut",
	//	})
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(resCr)
	//}
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
