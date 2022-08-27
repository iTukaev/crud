package helper

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/metadata"
)

const (
	undefinedMeta = "undefined"
)

func GetMetaFromContext(ctx context.Context) (context.Context, string) {
	var data []string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		data = md.Get("meta")
	}
	var meta string
	if len(data) > 0 {
		meta = data[0]
	} else {
		meta = fmt.Sprintf("%s_%d", undefinedMeta, time.Now().UTC().Unix())
	}

	return metadata.NewOutgoingContext(ctx, md), meta
}
