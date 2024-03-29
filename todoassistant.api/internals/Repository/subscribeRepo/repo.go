package subscribeRepo

import (
	"context"
	"test-va/internals/entity/subscribeEntity"
)

type SubscribeRepository interface {
	PersistEmail(ctx context.Context, req *subscribeEntity.SubscribeReq) error
	DeleteEmail(ctx context.Context, req *subscribeEntity.SubscribeReq) error
	CheckEmail(ctx context.Context, req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, error)
}
