package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/brokers/validator"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
)

func main() {
	config, err := yamlPkg.New()
	if err != nil {
		log.Fatalln("Config init error:", err)
	}
	logger, err := loggerPkg.New(config.LogLevel())
	if err != nil {
		log.Fatalln("Config init error:", err)
	}
	logger.Infoln("Start validator")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Infoln("Shutting down...")
		cancel()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err = runService(ctx, config, logger); err != nil {
			logger.Errorf("Consumer: %v", err)
		}
		c <- os.Interrupt
	}()

	<-c
}

func runService(ctx context.Context, config configPkg.Interface, logger *zap.SugaredLogger) error {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer(config.Brokers(), cfg)
	if err != nil {
		return errors.Wrap(err, "new SyncProducer")
	}

	income, err := sarama.NewConsumerGroup(config.Brokers(), consts.GroupValidator, cfg)
	if err != nil {
		return errors.Wrap(err, "new ConsumerGroup")
	}

	handler := validator.NewHandler(logger, producer)

	go func() {
		for {
			if err = income.Consume(ctx, []string{consts.TopicValidate}, handler); err != nil {
				logger.Errorf("on consume: <%v>", err)
				time.Sleep(time.Second * 5)
			}
		}
	}()

	<-ctx.Done()
	return income.Close()
}
