package main

import (
	"bot/config"
	"log"

	"bot/internal/commander"
	"bot/internal/handler"
	"bot/internal/storage"
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
