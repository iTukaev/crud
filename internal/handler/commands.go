package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"bot/internal/storage"
)

const (
	list   = "list"
	help   = "help"
	add    = "add"
	update = "update"
	del    = "del"
)

func listCommand(_ string) string {
	data := storage.List()
	res := make([]string, 0, len(data))
	for _, v := range data {
		res = append(res, v.String())
	}
	return strings.Join(res, "\n")
}

func helpCommand(_ string) string {
	return "/help - current menu\n" +
		"/list - all users\n" +
		"/add <name> <password> - add user\n" +
		"/update <id> <new name> <new password>\n" +
		"/del <id> - delete user"
}

func addCommand(data string) string {
	data = strings.TrimSpace(data)
	args := strings.Split(data, " ")
	if len(args) != 2 {
		return errors.Wrapf(BadCommand, "%d items: <%v>", len(args), args).Error()
	}

	u, err := storage.NewUser(args[0], args[1])
	if err != nil {
		return err.Error()
	}

	err = storage.Add(u)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("user: %s added", u.String())
}

func delCommand(data string) string {
	data = strings.TrimSpace(data)
	id, err := getIdFromString(data)
	if err != nil {
		return err.Error()
	}

	if u, err := storage.Delete(uint(id)); id < 0 || err != nil {
		return errors.Wrapf(err, "id <%d> invalid", id).Error()
	} else {
		return fmt.Sprintf("user: %s deleted", u.String())
	}
}

func updateCommand(data string) string {
	data = strings.TrimSpace(data)
	args := strings.Split(data, " ")
	if len(args) != 3 {
		return errors.Wrapf(BadCommand, "%d items: <%v>", len(args), args).Error()
	}

	id, err := getIdFromString(args[0])
	if err != nil {
		return err.Error()
	}

	if u, err := storage.Update(id, args[1], args[2]); err != nil {
		return errors.Wrap(err, "update error").Error()
	} else {
		return fmt.Sprintf("user: %s updated", u.String())
	}
}

func getIdFromString(data string) (uint, error) {
	id, err := strconv.Atoi(data)
	if err != nil {
		return 0, errors.Wrapf(err, "id <%s> is not a number", data)
	}
	if id < 0 {
		return 0, fmt.Errorf("id <%d> invalid, expected number > 0", id)
	}
	return uint(id), nil
}
