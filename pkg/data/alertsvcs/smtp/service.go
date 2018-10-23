package smtp

import (
	"crypto/tls"
	"errors"

	"gopkg.in/gomail.v2"
)

// SendEmail Sends Email
func SendEmail(c Config, subject string, body string) error {
	if !c.Enabled {
		return errors.New("SendEmail. Service is not enabled")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", c.To...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	//m.SetBody("text/plain", body)

	//m.SetAddressHeader("Cc", "peter@example.com", "Peter")
	//m.Attach("/home/Peter/attach.jpg")

	d := gomail.Dialer{Host: c.Host, Port: c.Port}
	if len(c.Username) > 0 {
		d.Username = c.Username
		d.Password = c.Password
	}

	d.TLSConfig = &tls.Config{InsecureSkipVerify: c.InsecureSkipVerify}

	// Send the email
	err := d.DialAndSend(m)

	if err != nil {
		return err
	}
	return nil
}
