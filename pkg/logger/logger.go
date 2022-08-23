package logger

type Interface interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}
