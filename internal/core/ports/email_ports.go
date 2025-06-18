package ports

import (
	"github.com/PatipanDev/mini-project-golang/internal/core/domain"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type EmailLogRepository interface {
	Save(log *domain.EmailLog) error
	ProcessUserEmailJob(args *domain.SendEmailArgs) error
}
