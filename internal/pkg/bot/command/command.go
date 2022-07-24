package command

type Interface interface {
	Process(args string) string
	Name() string
	Description() string
}
