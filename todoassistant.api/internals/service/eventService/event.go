package eventService

import (
	"context"
)

type TaskDistributor interface {
	SendTask(ctx context.Context, name string, payload []byte)
}

type TaskProcessor interface {
	ProcessTask(ctx context.Context, name string)
}
