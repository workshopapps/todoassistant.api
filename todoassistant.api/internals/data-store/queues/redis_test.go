package queues

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestNewRedisQueue(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 45*time.Second)
	defer cancelFunc()
	queue, err := NewRedisQueue("")
	if err != nil {
		log.Println("here")
		t.Fatal(err)
	}

	result, err := queue.Client.Ping(ctx).Result()
	if err != nil {
		t.Error("error", err)
		return
	}
	log.Println(result)
}
