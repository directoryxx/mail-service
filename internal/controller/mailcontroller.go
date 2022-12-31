package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"mail/internal/domain"
	"mail/internal/usecase"
)

// interface
type MailController interface {
	SendMail(ctx context.Context, data string) error
}

// implement interface
type MailControllerImpl struct {
	MailUsecase usecase.MailUsecase
}

func NewMailController(mailUsecase usecase.MailUsecase) MailController {
	return &MailControllerImpl{
		MailUsecase: mailUsecase,
	}
}

func (mc *MailControllerImpl) SendMail(ctx context.Context, data string) (err error) {
	messageDomain := &domain.Message{}
	var jsonData = []byte(data)

	var _ = json.Unmarshal(jsonData, &messageDomain)

	// Send Email
	err = mc.MailUsecase.SendMail(ctx, messageDomain)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
