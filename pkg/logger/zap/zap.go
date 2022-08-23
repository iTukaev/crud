package zap

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

var levelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func New(lvl string) (loggerPkg.Interface, error) {
	level := zap.NewAtomicLevel()
	level.SetLevel(getLoggerLevel(lvl))
	coreLog, err := zapcore.NewIncreaseLevelCore(zapcore.NewNopCore(), level)
	if err != nil {
		return nil, errors.Wrap(err, "increase logger level error")
	}
	logger := zap.New(coreLog)
	sugared := logger.Sugar()

	return &core{
		log: sugared,
	}, nil
}

type core struct {
	log *zap.SugaredLogger
}

func (c *core) Fatal(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c *core) Error(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c *core) Info(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (c *core) Debug(args ...interface{}) {
	//TODO implement me
	panic("implement me")
}
