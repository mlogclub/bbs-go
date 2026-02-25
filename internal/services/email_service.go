package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/repositories"
	"bytes"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/jordan-wright/email"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
)

var EmailService = newEmailService()

func newEmailService() *emailService {
	templateStr := `<div style="background-color:white;border-top:2px solid #12ADDB;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
    <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;">
        <span style="color: #12ADDB;font-weight:bold;">
            {{.Title}}
        </span>
    </h2>
    <div style="padding:0 12px 0 12px; margin-top:18px;">
		<p>
		{{if .From}}<a href="{{.From.Url}}" target="_blank" rel="noopener">{{.From.Title}}</a>&nbsp;{{end}}{{if .Content}}{{.Content}}{{end}}
		</p>
		{{if .QuoteContent}}
		<div style="background-color: #f5f5f5;padding: 10px 15px;margin:18px 0;word-wrap:break-word;">
            {{.QuoteContent}}
        </div>
		{{end}}
       
		{{if .link}}
        <p>
            <a style="text-decoration:none; color:#12addb" href="{{.link.Url}}" target="_blank" rel="noopener">{{.link.Title}}</a>
        </p>
		{{end}}
    </div>
</div>`
	tpl, err := template.New("emailTemplate").Parse(templateStr)
	if err != nil {
		slog.Error("init email template error", slog.Any("error", err))
	}
	return &emailService{
		tpl: tpl,
	}
}

type emailService struct {
	tpl *template.Template
}

func (s *emailService) SendTemplateEmail(from *models.User, to, subject, title, content, quote string, link *dto.ActionLink, bizType string) error {
	var fromLink *dto.ActionLink
	if from != nil {
		fromLink = &dto.ActionLink{
			Title: from.Nickname,
			Url:   bbsurls.UserUrl(from.Id),
		}
	}
	var b bytes.Buffer
	if err := s.tpl.Execute(&b, map[string]interface{}{
		"From":         fromLink,
		"Title":        title,
		"Content":      content,
		"QuoteContent": quote,
		"link":         link,
	}); err != nil {
		return err
	}
	html := b.String()
	return s.sendEmail(to, subject, html, bizType)
}

func (s *emailService) sendEmail(to string, subject, html string, bizType string) error {
	smtpConfig, configErr := s.getSmtpConfig()
	var (
		host      = smtpConfig.Host
		port      = smtpConfig.Port
		username  = smtpConfig.Username
		password  = smtpConfig.Password
		ssl       = smtpConfig.SSL
		addr      = net.JoinHostPort(host, port)
		auth      = smtp.PlainAuth("", username, password, host)
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
	)

	e := email.NewEmail()
	e.From = username
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)

	sendErr := error(nil)
	if configErr == nil {
		if ssl {
			sendErr = e.SendWithTLS(addr, auth, tlsConfig)
		} else {
			sendErr = e.Send(addr, auth)
		}
	} else {
		sendErr = configErr
	}

	if bizType == "" {
		bizType = constants.EmailLogBizTypeUnknown
	}

	emailLog := &models.EmailLog{
		ToEmail:    to,
		Subject:    subject,
		Content:    html,
		BizType:    bizType,
		Status:     constants.EmailLogStatusSuccess,
		CreateTime: dates.NowTimestamp(),
	}
	if sendErr != nil {
		emailLog.Status = constants.EmailLogStatusFailed
		emailLog.ErrorMsg = sendErr.Error()
	}

	if err := repositories.EmailLogRepository.Create(sqls.DB(), emailLog); err != nil {
		slog.Error("Create email log failed", slog.Any("err", err))
	}

	if sendErr != nil {
		slog.Error("Send email failed", slog.Any("err", sendErr))
		return sendErr
	}
	return nil
}

func (s *emailService) getSmtpConfig() (dto.SmtpConfig, error) {
	raw := cache.SysConfigCache.GetStr(constants.SysConfigSmtpConfig)
	var smtpConfig dto.SmtpConfig
	if strings.TrimSpace(raw) == "" {
		return smtpConfig, errors.New("smtp config is empty")
	}
	if err := jsons.Parse(raw, &smtpConfig); err != nil {
		return dto.SmtpConfig{}, err
	}
	if strings.TrimSpace(smtpConfig.Host) == "" ||
		strings.TrimSpace(smtpConfig.Port) == "" ||
		strings.TrimSpace(smtpConfig.Username) == "" ||
		strings.TrimSpace(smtpConfig.Password) == "" {
		return dto.SmtpConfig{}, errors.New("smtp config is incomplete")
	}
	return smtpConfig, nil
}
