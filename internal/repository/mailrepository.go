package repository

import (
	"bytes"
	"context"
	"fmt"
	"mail/internal/domain"
	"text/template"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository represent the user's repository contract
type MailRepository interface {
	SendMail(ctx context.Context, msg *domain.Message) error
	InsertMailLog(ctx context.Context, mailLog *domain.MailLog) error
}

type MailRepositoryImpl struct {
	Mail *mail.SMTPClient
	DB   *mongo.Client
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMailRepository(mail *mail.SMTPClient, db *mongo.Client) MailRepository {
	return &MailRepositoryImpl{
		Mail: mail,
		DB:   db,
	}
}

func (mailRepo *MailRepositoryImpl) InsertMailLog(ctx context.Context, log *domain.MailLog) (err error) {
	coll := mailRepo.DB.Database("logs").Collection("mail_log")
	_, err = coll.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	return nil

}

func (mailRepo *MailRepositoryImpl) SendMail(ctx context.Context, msg *domain.Message) (err error) {
	filter := bson.D{{"uuid", bson.D{{"$eq", msg.Uuid}}}}

	count, _ := mailRepo.DB.Database("logs").Collection("mail_log").CountDocuments(ctx, filter)

	// Return when it already inserted
	if count > 0 {
		return nil
	}

	// Insert Process Status
	mailLogProcess := &domain.MailLog{
		Uuid:   msg.Uuid,
		Status: "process",
	}

	mailRepo.InsertMailLog(ctx, mailLogProcess)

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := BuildHTMLMessage(msg)
	if err != nil {
		fmt.Println("Format Error :", err)
	}

	plainMessage, err := BuildPlainTextMessage(msg)
	if err != nil {
		fmt.Println("Plain Error :", err)
	}

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(mailRepo.Mail)
	if err != nil {
		return err
	}

	// Insert Send Status
	mailLogProcess = &domain.MailLog{
		Uuid:   msg.Uuid,
		Status: "send",
	}

	mailRepo.InsertMailLog(ctx, mailLogProcess)

	return nil

}

func BuildHTMLMessage(msg *domain.Message) (string, error) {
	templateToRender := "./internal/templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = InlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func BuildPlainTextMessage(msg *domain.Message) (string, error) {
	templateToRender := "./internal/templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func InlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}
