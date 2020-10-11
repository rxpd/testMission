package controllers

import (
	"avito/utils"
	"fmt"
	"net/smtp"
)

const (
	SMTPFrom     = "avito@293474-cd03243.tmweb.ru"
	SMTPPassword = "i41BEU8C"
	SMTPHost     = "smtp.timeweb.ru"
	SMTPPort     = "2525"
)

func SendConfirmationToken(email, token string) error {

	to := []string{
		email,
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Подтверждение email\n"
	preMessage := fmt.Sprintf("<html><body>Пожалуйста подтвердите ваш email, перейдя по <a href=\"127.0.0.1:8080/confirm_email?token=%s\" (%s)>этой ссылке</a></body></html>", token, token)
	message := []byte(subject + mime + preMessage)

	auth := smtp.PlainAuth("", SMTPFrom, SMTPPassword, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SMTPFrom, to, message)
	if err != nil {
		return err
	}
	utils.Log(preMessage)
	//fmt.Println("Email Sent Successfully!")
	return nil
}

func SendNotifyEmail(to []string, newPrice, oldPrice, url, title string) error {
	fmt.Println(to)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := fmt.Sprintf("Subject: Изменение цены на %s\n", title)
	preMessage := fmt.Sprintf("<html><body>Цена на <a href=\"%s\">%s</a> изменилась на с %s на %s ₽.</body></html>", url, title, oldPrice, newPrice)
	message := []byte(subject + mime + preMessage)

	auth := smtp.PlainAuth("", SMTPFrom, SMTPPassword, SMTPHost)
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, SMTPFrom, to, message)
	if err != nil {
		return err
	}
	utils.Log(preMessage)
	//fmt.Println("Email Sent Successfully!")
	return nil
}
