package main

import (
	"log"

	"gitlab.ozon.dev/kshmatov/masterclass1/internal/commander"
)

func main() {
	cmd, err := commander.Init()
	if err != nil {
		log.Panic(err)
	}
	if err := cmd.Run(); err != nil {
		log.Panic(err)
	}
}
