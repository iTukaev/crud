package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/kshmatov/masterclass1/internal/commander"

	"gitlab.ozon.dev/kshmatov/masterclass1/internal/storage"
)

const (
	helpCmd = "help"
	listCmd = "list"
	addCmd  = "add"
)

var BadArgument = errors.New("bad argument")

func listFunc(s string) string {
	data := storage.List()
	res := make([]string, 0, len(data))
	for _, v := range data {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}

func helpFunc(s string) string {
	return "/help - list commands\n" +
		"/list - list data\n" +
		"/add <name> <password> - add new user with name and password"
}

func addFunc(data string) string {
	log.Printf("add command param: <data>")
	params := strings.Split(data, " ")
	if len(params) != 2 {
		return errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params).Error()
	}
	u, err := storage.NewUser(params[0], params[1])
	if err != nil {
		return err.Error()
	}
	err = storage.Add(u)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("user %v added", u)
}

func AddHandlers(c *commander.Commander) {
	c.RegisterHandler(helpCmd, helpFunc)
	c.RegisterHandler(listCmd, listFunc)
	c.RegisterHandler(addCmd, addFunc)
}
