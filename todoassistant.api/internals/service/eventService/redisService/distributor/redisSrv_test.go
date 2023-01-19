package distributor

import (
	"log"
	"test-va/internals/data-store/queues"
	"testing"
	"time"
)

func Test_taskDistributor_SendTask(t *testing.T) {
	queue, err := queues.NewRedisQueue("")
	if err != nil {
		log.Println("here")
		t.Fatal(err)
	}

	distributor := NewTaskDistributor(queue.Client)

	for i := 0; i < 10; i++ {
		log.Println("sending Tasks")
		distributor.SendTask()
		time.Sleep(10 * time.Second)
	}
}
