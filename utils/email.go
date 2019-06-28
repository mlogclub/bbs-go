package utils

import (
	"github.com/jordan-wright/email"
	"github.com/mlogclub/mlog/utils/config"
	"net/smtp"
	"time"
)

var emailPool *email.Pool

func InitEmail() {
	addr := config.Conf.Smtp.Addr + ":" + config.Conf.Smtp.Port
	emailPool, _ = email.NewPool(addr, 3, smtp.PlainAuth("", config.Conf.Smtp.Username, config.Conf.Smtp.Password,
		config.Conf.Smtp.Addr))
}

func SendEmail(to string, subject, html string) error {
	e := email.NewEmail()
	e.From = config.Conf.Smtp.Username
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)
	return emailPool.Send(e, 10*time.Second)
}
