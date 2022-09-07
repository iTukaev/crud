package helper

import (
	"context"

	"github.com/Shopify/sarama"
)

const (
	uidKey = "uid"
	pubKey = "pub"
)

func InjectUidPubToCtx(ctx context.Context, uid, pub string) context.Context {
	ctx = context.WithValue(ctx, uidKey, uid)
	ctx = context.WithValue(ctx, pubKey, pub)
	return ctx
}

func ExtractUidPubFromCtx(ctx context.Context) (uid, pub string) {
	uid, _ = ctx.Value(uidKey).(string)
	pub, _ = ctx.Value(pubKey).(string)
	return uid, pub
}

func ExtractUidPubFromMessage(msg *sarama.ConsumerMessage) (uid, pub string) {
	for _, header := range msg.Headers {
		switch string(header.Key) {
		case uidKey:
			uid = string(header.Value)
		case pubKey:
			pub = string(header.Value)
		}
	}
	return uid, pub
}
