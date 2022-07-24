package help

import (
	"fmt"

	commandPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/bot/command"
)

func New(helpsMap map[string]string) commandPkg.Interface {
	if helpsMap == nil {
		helpsMap = make(map[string]string)
	}
	return &command{
		helps: helpsMap,
	}
}

type command struct {
	helps map[string]string
}

func (c *command) Process(_ string) string {
	result := fmt.Sprintf("/%s - %s", c.Name(), c.Description())
	for name, description := range c.helps {
		result += "\n"
		result += fmt.Sprintf("/%s - %s", name, description)
	}
	return result
}

func (*command) Name() string {
	return "help"
}

func (*command) Description() string {
	return "list of commands"
}
