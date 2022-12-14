package emailService

import (
	"net/smtp"
	"test-va/internals/entity/emailEntity"

)

type EmailService interface {
	SendMail(req emailEntity.SendEmailReq) error
	SendBatchEmail (req emailEntity.SendBatchEmail) error
}
type emailSrv struct {
	FromEmail string
	Password  string
	Host      string
	Port      string
}
func (e emailSrv) SendBatchEmail (req emailEntity.SendBatchEmail) error{
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	body := []byte(req.EmailSubject + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, req.EmailAddresses, body)
	if err != nil {
		return err
	}
	return nil
}
func (e emailSrv) SendMail(req emailEntity.SendEmailReq) error {
	auth := smtp.PlainAuth("", e.FromEmail, e.Password, e.Host)
	addr := e.Host + ":" + e.Port
	body := []byte(req.EmailSubject + req.EmailBody)
	err := smtp.SendMail(addr, auth, e.FromEmail, []string{req.EmailAddress}, body)
	if err != nil {
		return err
	}
	return nil
}

// func SendMail(req emailEntity.SendEmailReq) error {
// 	// add the email and password below, but remember to remove before pushing to giyhub
// 	from := utils.Config.FromEmailAddr
// 	password := os.Getenv("SMTPpwd")
// 	toEmail := req.EmailAddress
// 	to := []string{toEmail}
// 	host := "smtp.gmail.com"
// 	port := "587"
// 	address := host + ":" + port
// 	subject := req.EmailSubject
// 	body := req.EmailBody
// 	message := []byte(subject + "\n" + body)
// 	auth := smtp.PlainAuth("", from, password, host)
// 	err := smtp.SendMail(address, auth, from, to, message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func NewEmailSrv(fromEmail string, password string, host string, port string) EmailService {
	return &emailSrv{FromEmail: fromEmail, Password: password, Host: host, Port: port}
}
