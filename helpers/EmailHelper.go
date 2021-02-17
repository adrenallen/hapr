package helpers

import (
	gomail "gopkg.in/gomail.v2"
)

func NewEmailMessage(subject string, body string, toAddresses []string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "support@hapr.io", "Hapr")
	m.SetHeader("To", toAddresses...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return m
}

func NewReminderEmailMessage(subject string, body string, toAddresses []string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "reminders@hapr.io", "Hapr")
	m.SetAddressHeader("To", "reminders@hapr.io", "Hapr")
	m.SetHeader("Bcc", toAddresses...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return m
}

func SendEmailMessage(m *gomail.Message) error {
	emailConfig := getUserKeyConfig()
	d := gomail.NewDialer("smtp.gmail.com", 465, emailConfig.User, emailConfig.Key)

	return d.DialAndSend(m)
}

func getUserKeyConfig() EmailConfig {
	emailConfig := EmailConfig{}

	user, err := GetConfigValue("email.user")
	if err != nil {
		panic(err)
	}

	key, err := GetConfigValue("email.key")
	if err != nil {
		panic(err)
	}

	emailConfig.User = user
	emailConfig.Key = key
	return emailConfig

}

type EmailConfig struct {
	User string
	Key  string
}
