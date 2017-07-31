// Package mail is for removing the need to write boilerplate when sending email
// from a program within the Sanger network
package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
)

var (
	mailHost = "mail.sanger.ac.uk"
)

// Send provides a wrapper for the usual boilerplate to reduce the
// hassle of programatically sending emails within the Sanger.
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
