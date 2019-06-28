package oauth

import (
	"bytes"
	"github.com/mlogclub/mlog/utils/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"html/template"
)

func GetSuccessHtml(token *oauth2.Token, webUrl string) string {
	tplStr := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>登录成功</title>
    <link href="https://cdn.bootcss.com/bulma/0.7.2/css/bulma.min.css" rel="stylesheet">
</head>
<body>
<section class="section">
    <div class="container">
        <div class="notification is-primary">
            <strong>登录成功！</strong>
        </div>
    </div>
</section>
</body>
<script type="text/javascript">
  setTimeout(function () {
    window.location = '{{.WebUrl}}?tokenType={{.TokenType}}'
      + '&accessToken={{.AccessToken}}'
      + '&refreshToken={{.RefreshToken}}'
      + '&expiry={{.Expiry}}'
  }, 1000)
</script>
</html>
`
	tpl, err := template.New("successHtml").Parse(tplStr)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"WebUrl":       webUrl,
		"TokenType":    token.TokenType,
		"AccessToken":  token.AccessToken,
		"RefreshToken": token.RefreshToken,
		"Expiry":       token.Expiry,
	})
	if err != nil {
		logrus.Error(err)
		return ""
	}

	return b.String()
}

func GetOauthConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     config.Conf.OauthClient.ClientId,
		ClientSecret: config.Conf.OauthClient.ClientSecret,
		RedirectURL:  config.Conf.OauthClient.ClientRedirectUrl,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Conf.OauthServer.AuthUrl,
			TokenURL: config.Conf.OauthServer.TokenUrl,
		},
	}
	return config
}
