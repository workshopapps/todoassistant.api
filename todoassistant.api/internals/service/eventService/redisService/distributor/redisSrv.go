package distributor

import (
	"context"
	"log"
	"test-va/internals/service/eventService"

	"github.com/go-redis/redis/v9"
)

const TaskSendEmail = "send_email"

type taskDistributor struct {
	rdb *redis.Client
}

func (t taskDistributor) SendTask(ctx context.Context, name string, payload []byte) {
	publish := t.rdb.Publish(context.Background(), name, payload)
	err := publish.Err()
	if err != nil {
		log.Println(err)
		return
	}
}

func NewTaskDistributor(redisClient *redis.Client) eventService.TaskDistributor {
	return &taskDistributor{rdb: redisClient}
}
