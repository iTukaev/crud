package helper

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func StartNewSpan(ctx context.Context, name string, stop chan struct{}) {
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		pCtx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			span := tracer.StartSpan(name, opentracing.ChildOf(pCtx))
			defer func() {
				span.Finish()
			}()
		}
	}
	<-stop
}

func InjectHeaders(ctx context.Context, msg *sarama.ProducerMessage) error {
	span := opentracing.SpanFromContext(ctx)
	uid, pub := ExtractUidPubFromCtx(ctx)
	headers := map[string]string{
		uidKey: uid,
		pubKey: pub,
	}

	if err := opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.TextMap,
		opentracing.TextMapCarrier(headers),
	); err != nil {
		return errors.Wrap(err, "inject span to global tracer")
	}

	for key, value := range headers {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(value),
		})
	}
	return nil
}

func GetSpanFromMessage(msg *sarama.ConsumerMessage, operationName string) opentracing.Span {
	headers := make(map[string]string)
	for _, header := range msg.Headers {
		headers[string(header.Key)] = string(header.Value)
	}

	spanContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(headers))
	if err != nil {
		fmt.Println(err)
		return opentracing.StartSpan(operationName)
	}
	return opentracing.StartSpan(
		operationName,
		opentracing.FollowsFrom(spanContext),
	)
}
