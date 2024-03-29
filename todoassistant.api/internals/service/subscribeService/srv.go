package subscribeService

import (
	"context"
	"fmt"
	"log"
	"test-va/internals/Repository/subscribeRepo"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/eventEntity"
	"test-va/internals/entity/subscribeEntity"
	"test-va/internals/msg-queue/Emitter"
	"test-va/internals/service/emailService"
	"time"
)

type SubscribeService interface {
	PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError)
	DeleteSubscriber(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError)
	Contact(req *subscribeEntity.ContactUsReq) (*subscribeEntity.ContactUsRes, *ResponseEntity.ServiceError)
}

type subscribeSrv struct {
	repo     subscribeRepo.SubscribeRepository
	emailSrv emailService.EmailService
	Emitter  Emitter.Emitter
}

func NewSubscribeSrv(repo subscribeRepo.SubscribeRepository, emailSrv emailService.EmailService, emitter Emitter.Emitter) SubscribeService {
	return &subscribeSrv{repo: repo, emailSrv: emailSrv, Emitter: emitter}
}

// Subscribe to service godoc
// @Summary	Provide email to be subscribed to our service
// @Description	Add a subscriber route
// @Tags	Subscribe
// @Accept	json
// @Produce	json
// @Param	request	body	subscribeEntity.SubscribeReq	true	"Subscribe request"
// @Success	200  {object}  subscribeEntity.SubscribeRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/subscribe [post]
func (t *subscribeSrv) PersistEmail(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError) {
	// var message emailEntity.SendEmailReq
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	result, err1 := t.repo.CheckEmail(ctx, req)

	if result != nil {
		return nil, ResponseEntity.NewCustomServiceError("Already subscribed", err1)
	}

	// message.EmailAddress = req.Email
	// message.EmailSubject = "Subject: Subscription To Ticked Newsletter\n"
	// message.EmailBody = CreateMessageBody()

	//err := t.emailSrv.SendMail(message)
	//if err != nil {
	//	return nil, ResponseEntity.NewInternalServiceError(err)
	//}

	err := t.repo.PersistEmail(ctx, req)
	if err != nil {
		log.Println("From subscribe ", err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := subscribeEntity.SubscribeRes{
		Email: req.Email,
	}

	// push event to queue
	payload := eventEntity.Payload{
		Action:    "email",
		SubAction: "subscription",
		Data: map[string]string{
			"email_address": req.Email,
			"email_subject": "Subject: Subscription To Ticked Newsletter\n",
			"email_body":    CreateMessageBody(),
		},
	}

	err = t.Emitter.Push(payload, "info")
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return &data, nil
}

// Subscribe to service godoc
// @Summary	Provide email to be subscribed to our service
// @Description	Add a subscriber route
// @Tags	Subscribe
// @Accept	json
// @Produce	json
// @Param	request	body	subscribeEntity.SubscribeReq	true	"Subscribe request"
// @Success	200  {object}  subscribeEntity.SubscribeRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/subscribe [delete]
func (t *subscribeSrv) DeleteSubscriber(req *subscribeEntity.SubscribeReq) (*subscribeEntity.SubscribeRes, *ResponseEntity.ServiceError) {
	// var message emailEntity.SendEmailReq
	// create context of 1 minute
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	result, err1 := t.repo.CheckEmail(ctx, req)

	if result == nil {
		return nil, ResponseEntity.NewCustomServiceError("No subscriber with that email address", err1)
	}

	err := t.repo.DeleteEmail(ctx, req)
	if err != nil {
		log.Println("From subscribe ", err)
		return nil, ResponseEntity.NewInternalServiceError(err)
	}
	data := subscribeEntity.SubscribeRes{
		Email: req.Email,
	}

	// push event to queue
	payload := eventEntity.Payload{
		Action:    "email",
		SubAction: "subscription",
		Data: map[string]string{
			"email_address": req.Email,
			"email_subject": "Subject: Unsubscription To Ticked Newsletter\n",
			"email_body":    CreateMessageBodyUnsubscribe(),
		},
	}

	err = t.Emitter.Push(payload, "info")
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return &data, nil
}

// Contact us godoc
// @Summary	Contact us with any complaint or feedback
// @Description	Contact us route
// @Tags	Contact-Us
// @Accept	json
// @Produce	json
// @Param	request	body	subscribeEntity.ContactUsReq	true	"Contact Us request"
// @Success	200  {object}  subscribeEntity.ContactUsRes
// @Failure	400  {object}  ResponseEntity.ServiceError
// @Failure	404  {object}  ResponseEntity.ServiceError
// @Failure	500  {object}  ResponseEntity.ServiceError
// @Router	/contact-us [post]
func (t *subscribeSrv) Contact(req *subscribeEntity.ContactUsReq) (*subscribeEntity.ContactUsRes, *ResponseEntity.ServiceError) {
	// create context of 1 minute
	_, cancelFunc := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancelFunc()

	// err := t.emailSrv.SendMailToSupport(message)
	// if err != nil {
	// 	return nil, ResponseEntity.NewInternalServiceError(err)
	// }

	data := subscribeEntity.ContactUsRes{
		Email:   req.Email,
		Name:    req.Name,
		Message: req.Message,
	}

	// push event to queue
	payload := eventEntity.Payload{
		Action:    "email",
		SubAction: "contact-us",
		Data: map[string]string{
			"email_address": "admin@getticked.com",
			"email_subject": fmt.Sprintf("Subject: %s (%s)\n", req.Name, req.Email),
			"email_body":    CreateMessageBodySupport(req),
		},
	}

	err := t.Emitter.Push(payload, "info")
	if err != nil {
		return nil, ResponseEntity.NewInternalServiceError(err)
	}

	return &data, nil
}

// Auxillary function
func CreateMessageBody() string {
	subject := "Subscription to Ticked!\n\n"
	mainBody := "Thank you for subscribing to our newsletter!\n\nGet ready for an awesome ride"
	message := subject + mainBody
	return string(message)
}

func CreateMessageBodyUnsubscribe() string {
	subject := "Unsubscription to Ticked!\n\n"
	mainBody := "You have successfully unsubscribed to ticked newsletter, it was fun having you with us and we will miss you!\n\nThanks from all of us at Ticked!"
	message := subject + mainBody
	return string(message)
}

func CreateMessageBodySupport(req *subscribeEntity.ContactUsReq) string {
	subject := fmt.Sprintf("Subject: %s (%s)\n", req.Name, req.Email)
	// title := "" // fmt.Sprintf("Feedback from %s\n\n", req.Email)
	message := subject + req.Message
	return string(message)
}
