package logger

type Interface interface {
	Error(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Errorf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Close()
}
