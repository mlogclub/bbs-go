package email

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/common/config"
)

var emailTemplate = `
<div style="background-color:white;border-top:2px solid #12ADDB;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
    <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;">
        <span style="color: #12ADDB;font-weight:bold;">
            {{.Title}}
        </span>
    </h2>
    <div style="padding:0 12px 0 12px; margin-top:18px;">
        {{if .Content}}
		<p>
            {{.Content}}
        </p>
		{{end}}
		{{if .QuoteContent}}
		<div style="background-color: #f5f5f5;padding: 10px 15px;margin:18px 0;word-wrap:break-word;">
            {{.QuoteContent}}
        </div>
		{{end}}
       
		{{if .Url}}
        <p>
            <a style="text-decoration:none; color:#12addb" href="{{.Url}}" target="_blank" rel="noopener">点击查看详情</a>
        </p>
		{{end}}
    </div>
</div>
`

// 发送模版邮件
func SendTemplateEmail(to, subject, title, content, quoteContent, url string) {
	tpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		logrus.Error(err)
		return
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"Title":        title,
		"Content":      content,
		"QuoteContent": quoteContent,
		"Url":          url,
	})
	if err != nil {
		logrus.Error(err)
		return
	}
	html := b.String()
	SendEmail(to, subject, html)
}

// 发送邮件
func SendEmail(to string, subject, html string) {
	var (
		host      = config.Conf.Smtp.Host
		port      = config.Conf.Smtp.Port
		username  = config.Conf.Smtp.Username
		password  = config.Conf.Smtp.Password
		ssl       = config.Conf.Smtp.SSL
		addr      = net.JoinHostPort(host, port)
		auth      = smtp.PlainAuth("", username, password, host)
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
	)

	e := email.NewEmail()
	e.From = config.Conf.Smtp.Username
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)

	if ssl {
		if err := e.SendWithTLS(addr, auth, tlsConfig); err != nil {
			logrus.Error("发送邮件异常", err)
		}
	} else {
		if err := e.Send(addr, auth); err != nil {
			logrus.Error("发送邮件异常", err)
		}
	}
}
