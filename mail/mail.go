// Package mail is for removing the need to write boilerplate when sending email from a CLI app
// This implementation only allows for anonymous connection to a mailserver and is of limited value
// I will probably improve it at some point but for now it serves my purposes.
package mail

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

var (
	mailHost = ""
	// ErrNoMailHost gets returned when the mailHost has not been set
	ErrNoMailHost = errors.New("No mail host set. Use mail.SetHost() to assign a mail host")
)

// SetHost allows the developer to set the mail host through which sent mail is routed
func SetHost(str string) {
	mailHost = str
}

// Send provides a wrapper for the usual boilerplate to reduce the
// hassle of programatically sending emails. Don't forget to assign
// a mailhost BEFORE attempting to send mail
func Send(from, to, subject, msg string, cc []string) error {
	if mailHost == "" {
		return ErrNoMailHost
	}

	// Format the msg text so we get a subject
	msgFmt := fmt.Sprintf("To: %s\r\nCC: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, strings.Join(cc, ","), subject, msg)

	// Connect to mail host
	c, err := smtp.Dial(fmt.Sprintf("%s:25", mailHost))
	if err != nil {
		return err
	}
	defer c.Close()

	// Set sender and recipient
	c.Mail(from)
	c.Rcpt(to)

	// Stream the body
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	buf := bytes.NewBufferString(msgFmt)
	if _, err := buf.WriteTo(wc); err != nil {
		return err
	}

	return nil
}
