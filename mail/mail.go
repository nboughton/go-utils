// Package mail is for removing the need to write boilerplate when sending email from a CLI app
package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
)

var (
	mailHost = ""
)

// SetHost allows the developer to set the mail host through which sent mail is routed
func SetHost(str string) {
	mailHost = str
}

// Send provides a wrapper for the usual boilerplate to reduce the
// hassle of programatically sending emails. Don't forget to assign
// a mailhost BEFORE attempting to send mail
func Send(from, to, subject, msg string) error {
	// Format the msg text so we get a subject
	msgFmt := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, msg)

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
