package smtp

import (
	"net/smtp"
	"os"

	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
	"github.com/PatipanDev/mini-project-golang/internal/core/ports"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// adapters/smtp/smtp_sender.go
type SMTPSender struct {
	emailLogRepo ports.EmailLogRepository
	from         string
	addr         string
	auth         smtp.Auth
}

func NewSMTPSender(emailLogRepo ports.EmailLogRepository) *SMTPSender {
	_ = godotenv.Load()

	from := os.Getenv("SMTP_SENDER")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	addr := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	return &SMTPSender{
		emailLogRepo: emailLogRepo,
		from:         from,
		addr:         addr,
		auth:         auth,
	}
}

func (s *SMTPSender) SendEmail(to, subject, htmlBody string) error {
	msg := []byte("From: " + s.from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		htmlBody + "\r\n")

	err := smtp.SendMail(s.addr, s.auth, s.from, []string{to}, msg)

	// Always log, even if error
	_ = s.emailLogRepo.Save(&domain.EmailLog{
		EventID: uuid.New(),
		Subject: subject,
		From:    s.from,
		To:      to,
		Body:    htmlBody,
	})

	return err
}
