package utils

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
)

var emailPool *email.Pool

func initEmail() {
	addr := Conf.SmtpAddr + ":" + Conf.SmtpPort
	emailPool, _ = email.NewPool(addr, 3, smtp.PlainAuth("", Conf.SmtpUsername, Conf.SmtpPassword, Conf.SmtpAddr))
}

func SendEmail(to string, subject, html string) error {
	e := email.NewEmail()
	e.From = Conf.SmtpUsername
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)
	return emailPool.Send(e, 10*time.Second)
}
