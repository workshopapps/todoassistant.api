package mail_grpc

import (
	"context"
	"emailSrv/internal/entities/mailEntity"
	grpcMail "emailSrv/internal/grpc-mail"
	mailSerivce "emailSrv/internal/service/mail-serivce"
	"fmt"
	"log"
)

type MailServer struct {
	grpcMail.UnimplementedMailServiceServer
	mailSrv mailSerivce.EmailService
}

func NewMailGrpcServer(mailSrv mailSerivce.EmailService) *MailServer {
	return &MailServer{mailSrv: mailSrv}
}

func (m MailServer) SendBatchEmails(ctx context.Context, request *grpcMail.BatchMailRequest) (*grpcMail.MailResponse, error) {
	input := request.BatchMailEntry
	//log.Println(input, "input")
	mailEntry := mailEntity.SendBatchEmail{
		Name:           input.Name,
		EmailAddresses: input.EmailAddresses,
		EmailSubject:   input.EmailSubject,
		EmailBody:      input.EmailBody,
	}
	err := m.mailSrv.SendBatchEmail(mailEntry)
	if err != nil {
		res := &grpcMail.MailResponse{Result: "Failed to send batch email"}
		return res, err
	}
	return &grpcMail.MailResponse{Result: "batch Email sent successfully"}, nil
}

func (m MailServer) SendEmail(ctx context.Context, request *grpcMail.MailRequest) (*grpcMail.MailResponse, error) {
	input := request.GetMailEntry()
	mailEntry := mailEntity.SendEmailReq{
		Name:         input.Name,
		EmailAddress: input.EmailAddress,
		EmailSubject: input.EmailSubject,
		EmailBody:    input.EmailBody,
	}

	err := m.mailSrv.SendMail(mailEntry)
	if err != nil {
		log.Println("here")
		errorStmt := fmt.Sprintf("cannot send email: %w", err)
		res := &grpcMail.MailResponse{Result: errorStmt}
		return res, err
	}

	return &grpcMail.MailResponse{Result: "Email sent successfully"}, nil
}
