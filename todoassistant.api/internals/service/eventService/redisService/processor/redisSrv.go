package processor

import (
	"context"
	"log"
	"test-va/internals/entity/emailEntity"
	"test-va/internals/service/emailService"
	"test-va/internals/service/eventService"
)

type taskProcessor struct {
	//server  *asynq.Server
	rdb     *redis.Client
	mailSrv emailService.EmailService
}

func (t *taskProcessor) ProcessTask(ctx context.Context, name string) {
	pubsub := t.rdb.Subscribe(ctx, name)
	ch := pubsub.Channel()
	for msg := range ch {
		log.Println("message is: ", msg)
		log.Println("message payload  is: ", msg.Payload)
		log.Println("message is channel is: ", msg.Channel)

		switch msg.Channel {
		case "email":
			var email emailEntity.SendEmailReq
			//send email now
			t.mailSrv.SendMail(email)
		}

	}

}

func NewTaskProcessor(redisClient *redis.Client) eventService.TaskProcessor {
	return &taskProcessor{rdb: redisClient}
}
