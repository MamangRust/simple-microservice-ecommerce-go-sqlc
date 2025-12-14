package mailer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/smtp"
)

type Mailer struct {
	ctx      context.Context
	server   string
	port     int
	user     string
	password string
}

func NewMailer(ctx context.Context, server string, port int, user string, password string) *Mailer {
	return &Mailer{
		ctx:      ctx,
		server:   server,
		port:     port,
		user:     user,
		password: password,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	log.Printf("[Mailer] Sending email to %s via %s:%d", to, m.server, m.port)

	auth := smtp.PlainAuth("", m.user, m.password, m.server)

	headers := map[string]string{
		"From":         m.user,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/html; charset="UTF-8"`,
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", m.server, m.port)

	if err := smtp.SendMail(addr, auth, m.user, []string{to}, msg.Bytes()); err != nil {
		log.Printf("[Mailer] ❌ Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("[Mailer] ✅ Email sent successfully to %s", to)
	return nil
}
