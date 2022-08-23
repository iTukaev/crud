package emptylog

import (
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

func New() loggerPkg.Interface {
	return &core{}
}

type core struct{}

func (*core) Fatal(_ ...interface{}) {}

func (*core) Error(_ ...interface{}) {}

func (*core) Info(_ ...interface{}) {}

func (*core) Debug(_ ...interface{}) {}
