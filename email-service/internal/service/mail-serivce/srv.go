package mail_serivce

import (
	"emailSrv/internal/entities/mailEntity"
	"fmt"
	"log"
	"net/smtp"
)

type EmailService interface {
	SendMail(req mailEntity.SendEmailReq) error
	SendBatchEmail(req mailEntity.SendBatchEmail) error
}
type emailSrv struct {
	FromEmail string
	Password  string
	Host      string
	Port      string
}

func (e emailSrv) SendBatchEmail(req mailEntity.SendBatchEmail) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	header := fmt.Sprintf("From: %v\nTo: %v\n", e.FromEmail, req.EmailAddresses)
	body := []byte(header + req.EmailSubject + "\n" + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, req.EmailAddresses, body)
	if err != nil {
		log.Println("error", err)
		return err
	}
	return nil
}
func (e emailSrv) SendMail(req mailEntity.SendEmailReq) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	header := fmt.Sprintf("From: %v\nTo: %v\n", e.FromEmail, req.EmailAddress)
	body := []byte(header + req.EmailSubject + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, []string{req.EmailAddress}, body)
	if err != nil {
		log.Println("error", err)
		return err
	}
	return nil
}

func NewEmailSrv(fromEmail string, password string, host string, port string) EmailService {
	return &emailSrv{FromEmail: fromEmail, Password: password, Host: host, Port: port}
}
