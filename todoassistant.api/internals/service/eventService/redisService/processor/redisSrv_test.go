package processor

import (
	"context"
	"log"
	"test-va/internals/data-store/queues"
	"testing"
)

func Test_taskProcessor_ProcessTask(t *testing.T) {
	queue, err := queues.NewRedisQueue("")
	if err != nil {
		log.Println("here")
		t.Fatal(err)
	}
	processor := NewTaskProcessor(queue.Client)
	log.Println("recieving Tasks")
	processor.ProcessTask(context.Background(), "James")
}
