// Package email sends transactional mail over SMTP. In dev/CI this targets
// Mailpit (docker-compose, no auth); prod points SMTP_ADDR at a real relay.
// No dependency — net/smtp is enough for the one message we send.
package email

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

// SMTPMailer sends mail via a plain SMTP server. The zero auth (nil) matches
// Mailpit and most internal relays; add smtp.Auth here when a relay needs it.
type SMTPMailer struct {
	Addr string // host:port
	From string // From header
}

// SendVerification emails a one-click verification link. Kept plain-text: the
// link is the whole payload, and plain text dodges HTML-rendering surprises.
func (m SMTPMailer) SendVerification(to, link string) error {
	// The SMTP envelope sender (MAIL FROM) must be a bare address; only the
	// From header may carry a display name like "Bibseller <noreply@…>".
	envelopeFrom := m.From
	if addr, err := mail.ParseAddress(m.From); err == nil {
		envelopeFrom = addr.Address
	}
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Verify your Bibseller email\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n\r\n"+
			"Welcome to Bibseller. Confirm your email address:\r\n\r\n%s\r\n\r\n"+
			"This link expires in 24 hours. If you didn't sign up, ignore this email.\r\n",
		m.From, to, link)
	return smtp.SendMail(m.Addr, nil, envelopeFrom, []string{to}, []byte(msg))
}
