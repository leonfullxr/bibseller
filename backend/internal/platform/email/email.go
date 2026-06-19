// Package email sends transactional mail over SMTP. In dev/CI this targets
// Mailpit (docker-compose, no auth); prod points SMTP_ADDR at a real relay.
// No dependency - net/smtp is enough for the messages we send.
//
// Copy is keyed by recipient locale (D17). M8.1 ships English only; Spanish is
// added per-map in M8.2, and any locale without an entry falls back to English.
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

// content is one localized message: a subject and a plain-text body whose single
// %s is the link (the whole payload of every message we send).
type content struct {
	subject string
	body    string // exactly one %s, for the link
}

// pick returns the locale's copy, falling back to English for any unmapped one.
func pick(byLocale map[string]content, locale string) content {
	if c, ok := byLocale[locale]; ok {
		return c
	}
	return byLocale["en"]
}

var verificationEmail = map[string]content{
	"en": {
		subject: "Verify your Bibseller email",
		body: "Welcome to Bibseller. Confirm your email address:\r\n\r\n%s\r\n\r\n" +
			"This link expires in 24 hours. If you didn't sign up, ignore this email.\r\n",
	},
}

var passwordResetEmail = map[string]content{
	"en": {
		subject: "Reset your Bibseller password",
		body: "Someone asked to reset your Bibseller password. Set a new one:\r\n\r\n%s\r\n\r\n" +
			"This link expires in 1 hour. If you didn't ask for this, ignore this email.\r\n",
	},
}

var newMessageEmail = map[string]content{
	"en": {
		subject: "New message about your Bibseller listing",
		body:    "A buyer started a conversation about one of your listings. Read and reply:\r\n\r\n%s\r\n",
	},
}

// SendVerification emails a one-click verification link in the recipient's locale.
func (m SMTPMailer) SendVerification(to, link, locale string) error {
	return m.send(to, pick(verificationEmail, locale), link)
}

// SendPasswordReset emails a one-click password-reset link in the recipient's locale.
func (m SMTPMailer) SendPasswordReset(to, link, locale string) error {
	return m.send(to, pick(passwordResetEmail, locale), link)
}

// SendNewMessage notifies a seller, in their locale, that a buyer started a
// conversation about one of their listings; the inbox link is the payload.
func (m SMTPMailer) SendNewMessage(to, link, locale string) error {
	return m.send(to, pick(newMessageEmail, locale), link)
}

// send builds the plain-text envelope and delivers one message. The SMTP
// envelope sender (MAIL FROM) must be a bare address; only the From header may
// carry a display name like "Bibseller <noreply@…>". Plain text dodges
// HTML-rendering surprises - the link is the whole payload.
func (m SMTPMailer) send(to string, c content, link string) error {
	envelopeFrom := m.From
	if addr, err := mail.ParseAddress(m.From); err == nil {
		envelopeFrom = addr.Address
	}
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n"+
			"Content-Type: text/plain; charset=utf-8\r\n\r\n%s",
		m.From, to, c.subject, fmt.Sprintf(c.body, link))
	return smtp.SendMail(m.Addr, nil, envelopeFrom, []string{to}, []byte(msg))
}
