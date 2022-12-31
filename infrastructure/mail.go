package infrastructure

import (
	"crypto/tls"
	"os"
	"strconv"

	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func ConnectSMTP() (client *mail.SMTPClient, err error) {
	server := mail.NewSMTPClient()

	mailPort, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	// SMTP Server
	server.Host = os.Getenv("MAIL_SERVER")
	server.Port = mailPort
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionNone

	// Variable to keep alive connection
	server.KeepAlive = true

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	return smtpClient, err
}
