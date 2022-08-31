package counter

import (
	"fmt"
	"strings"
	"sync"
)

type core struct {
	mu   sync.Mutex
	data map[string]uint64
}

var (
	Request  *core
	Response *core
	Success  *core
	Errors   *core
)

func init() {
	Request = new(core)
	Request.data = make(map[string]uint64)

	Response = new(core)
	Response.data = make(map[string]uint64)

	Success = new(core)
	Success.data = make(map[string]uint64)

	Errors = new(core)
	Errors.data = make(map[string]uint64)
}

func (c *core) Inc(param string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[param]++
}

func (c *core) String() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := make([]string, 0, len(c.data))
	for key, val := range c.data {
		data = append(data, fmt.Sprintf("[%s] = %d", key, val))
	}
	return strings.Join(data, " | ")
}
