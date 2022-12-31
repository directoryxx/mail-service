package usecase

import (
	"context"
	"mail/internal/domain"
	"mail/internal/repository"
)

// UserRepository represent the user's repository contract
type MailUsecase interface {
	SendMail(ctx context.Context, msg *domain.Message) error
}

type MailUsecaseImpl struct {
	MailRepo repository.MailRepository
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMailUseCase(MailRepo repository.MailRepository) MailUsecase {
	return &MailUsecaseImpl{
		MailRepo: MailRepo,
	}
}

func (muc *MailUsecaseImpl) SendMail(ctx context.Context, msg *domain.Message) error {
	err := muc.MailRepo.SendMail(ctx, msg)
	return err
}
