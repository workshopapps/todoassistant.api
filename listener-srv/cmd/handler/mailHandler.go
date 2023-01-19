package handler

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"listener-srv/internal/entity/eventEntity"
	"listener-srv/internal/entity/mailEntity"
	grpc_mail "listener-srv/internal/grpc-mail"
	"log"
	"time"
)

type MailHandler struct {
	conn *grpc.ClientConn
}

func NewMailHandler() (*MailHandler, error) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*30)
	defer cancelFunc()
	conn, err := grpc.DialContext(ctx, "localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &MailHandler{conn}, nil
}

func (m *MailHandler) SendMail(payload eventEntity.Payload) (*grpc_mail.MailResponse, error) {
	var emailReq mailEntity.SendEmailReq

	emailReq.EmailAddress = payload.Data["email_address"]
	emailReq.EmailSubject = payload.Data["email_subject"]
	emailReq.EmailBody = payload.Data["email_body"]
	emailReq.Name = ""

	log.Printf("%#v", emailReq)

	c := grpc_mail.NewMailServiceClient(m.conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*30)
	defer cancelFunc()

	data, err := c.SendEmail(ctx, &grpc_mail.MailRequest{
		MailEntry: &grpc_mail.Mail{
			Name:         emailReq.Name,
			EmailAddress: emailReq.EmailAddress,
			EmailSubject: emailReq.EmailSubject,
			EmailBody:    emailReq.EmailBody,
		},
	})
	if err != nil {
		log.Println(data)
		return data, err
	}

	return data, nil
}

func (m *MailHandler) CloseConn() {
	m.conn.Close()

}
