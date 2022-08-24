package zaplog

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
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(getLoggerLevel(lvl)),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "lvl",
			TimeKey:        "time",
			NameKey:        "log",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build new logger")
	}
	sugared := logger.Sugar()

	return &core{
		log: sugared,
	}, nil
}

type core struct {
	log *zap.SugaredLogger
}

func (c *core) Error(args ...interface{}) {
	c.log.Errorln(args)
}

func (c *core) Info(args ...interface{}) {
	c.log.Infoln(args)
}

func (c *core) Debug(args ...interface{}) {
	c.log.Debugln(args)
}

func (c *core) Errorf(template string, args ...interface{}) {
	c.log.Errorf(template, args)
}

func (c *core) Infof(template string, args ...interface{}) {
	c.log.Infof(template, args)
}

func (c *core) Debugf(template string, args ...interface{}) {
	c.log.Debugf(template, args)
}

func (c *core) Close() {
	_ = c.log.Sync()
}
