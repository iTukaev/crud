package command

import "context"

type Interface interface {
	Process(ctx context.Context, args string) string
	Name() string
	Description() string
}
