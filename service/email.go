package service

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"server-article/model"
	"server-article/utils"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

type EmailDataResetPassword struct {
	URL     string
	Code    string
	Subject string
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *model.User, data *EmailData) {
	// config email
	from := utils.GetEnv("EMAIL_FROM")
	smtpHost := utils.GetEnv("EMAIL_HOST")
	smtpPort := utils.GetEnv("EMAIL_PORT")
	smtpUser := utils.GetEnv("EMAIL_USERNAME")
	smtpPass := utils.GetEnv("EMAIL_PASS")
	to := user.Email

	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("could not parse template", err)
		return
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	// str to int port
	port := utils.Str2Int(smtpPort)

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("verifivation_account: could not send email", err)
	}

}

func SendEmailResetPassword(user *model.User, data *EmailDataResetPassword) {
	// config email
	from := utils.GetEnv("EMAIL_FROM")
	smtpHost := utils.GetEnv("EMAIL_HOST")
	smtpPort := utils.GetEnv("EMAIL_PORT")
	smtpUser := utils.GetEnv("EMAIL_USERNAME")
	smtpPass := utils.GetEnv("EMAIL_PASS")
	to := user.Email

	var body bytes.Buffer

	// parse template
	template, err := ParseTemplateDir("templates_reset_pass")
	if err != nil {
		log.Fatal("could not parse template", err)
		return
	}

	template.ExecuteTemplate(&body, "resetPassword.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	// str to int port
	port := utils.Str2Int(smtpPort)

	d := gomail.NewDialer(smtpHost, port, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// send email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("reset_password: could not send email", err)
	}
}
