package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"gitlab.ozon.dev/iTukaev/homework/internal/brokers/mailing"
	configPkg "gitlab.ozon.dev/iTukaev/homework/internal/config"
	yamlPkg "gitlab.ozon.dev/iTukaev/homework/internal/config/yaml"
	"gitlab.ozon.dev/iTukaev/homework/internal/consts"
	jaegerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/jaeger"
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
	logger.Infoln("Start mailing")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Infoln("Shutting down...")
		cancel()
	}()

	tracer, closer, err := jaegerPkg.New(config.JService(), config.JHost())
	if err != nil {
		logger.Errorf("Jaeger initialise err: %v", err)
		return
	}
	defer func() {
		_ = closer.Close()
	}()
	opentracing.SetGlobalTracer(tracer)

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

	income, err := sarama.NewConsumerGroup(config.Brokers(), consts.GroupMailing, cfg)
	if err != nil {
		return errors.Wrap(err, "new ConsumerGroup")
	}

	handler := mailing.NewHandler(logger, producer)

	go func() {
		for {
			if err = income.Consume(ctx, []string{consts.TopicError, consts.TopicMailing}, handler); err != nil {
				logger.Errorf("on consume: <%v>", err)
				time.Sleep(time.Second * 5)
			}
		}
	}()

	<-ctx.Done()
	return income.Close()
}
