package main

import (
	"log"

	"gitlab.ozon.dev/iTukaev/homework-1/config"
	"gitlab.ozon.dev/iTukaev/homework-1/internal/commander"
	"gitlab.ozon.dev/iTukaev/homework-1/internal/handler"
	"gitlab.ozon.dev/iTukaev/homework-1/internal/storage"
)

func main() {
	log.Println("start main")

	if err := storage.Init(); err != nil {
		log.Panic(err)
	}

	if err := config.Init(); err != nil {
		log.Panic(err)
	}

	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	handler.AddHandlers(cmd)
	cmd.Run()
}
