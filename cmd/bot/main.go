package main

import (
	"log"

	"gitlab.ozon.dev/kshmatov/masterclass1/internal/commander"
	"gitlab.ozon.dev/kshmatov/masterclass1/internal/handlers"
)

func main() {
	log.Println("start main")
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	handlers.AddHandlers(cmd)

	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}
