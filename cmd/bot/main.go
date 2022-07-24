package main

import (
	"log"

	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	botPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/delete"
	cmdGetPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/get"
	cmdHelpPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command/update"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
)

func main() {
	log.Println("start main")
	user := userPkg.MustNew()

	config := yamlPkg.MustNew()
	config.Init()

	bot := botInit(user, config.BotKey())

	bot.Run()
}

func botInit(user userPkg.Interface, apiKey string) botPkg.Interface {
	bot := botPkg.MustNew(apiKey)

	commandAdd := cmdAddPkg.New(user)
	bot.RegisterCommander(commandAdd)

	commandUpdate := cmdUpdatePkg.New(user)
	bot.RegisterCommander(commandUpdate)

	commandDelete := cmdDeletePkg.New(user)
	bot.RegisterCommander(commandDelete)

	commandGet := cmdGetPkg.New(user)
	bot.RegisterCommander(commandGet)

	commandList := cmdListPkg.New(user)
	bot.RegisterCommander(commandList)

	commandHelp := cmdHelpPkg.New(map[string]string{
		commandAdd.Name():    commandAdd.Description(),
		commandUpdate.Name(): commandUpdate.Description(),
		commandDelete.Name(): commandDelete.Description(),
		commandGet.Name():    commandGet.Description(),
		commandList.Name():   commandList.Description(),
	})
	bot.RegisterCommander(commandHelp)

	return bot
}
