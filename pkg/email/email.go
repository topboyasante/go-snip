package email

import (
	"bytes"
	"fmt"
	"github.com/topboyasante/go-snip/pkg/config"
	"html/template"
	"log"
	"net/smtp"
)

type SMTPAuth struct {
	identity string
	username string
	password string
	host     string
}

var EmailConfig = initSMTPAuth()

func initSMTPAuth() *SMTPAuth {
	auth := &SMTPAuth{
		identity: "",
		username: config.ENV.SMTPUsername,
		password: config.ENV.SMTPPassword,
		host:     config.ENV.SMTPHost,
	}

	return auth
}

func SendMailWithSMTP(authVar *SMTPAuth, subject, templatePath string, values any, to []string) {
	var body bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("error: unable to parse template")
	}
	t.Execute(&body, values)

	auth := smtp.PlainAuth(
		authVar.identity,
		authVar.username,
		authVar.password,
		authVar.host,
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := fmt.Sprintf("Subject: %v\n%v \n\n%v", subject, headers, body.String())

	err = smtp.SendMail(
		config.ENV.SMTPAddress,
		auth,
		config.ENV.SMTPUsername,
		to,
		[]byte(msg),
	)

	if err != nil {
		fmt.Println(err)
	}
}
