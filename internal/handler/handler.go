package handler

import (
	"github.com/pkg/errors"

	"bot/internal/commander"
)

var BadCommand = errors.New("unexpected arguments")

func AddHandlers(c *commander.Commander) {
	c.RegisterCommander(list, listCommand)
	c.RegisterCommander(help, helpCommand)
	c.RegisterCommander(add, addCommand)
	c.RegisterCommander(del, delCommand)
	c.RegisterCommander(update, updateCommand)
}
