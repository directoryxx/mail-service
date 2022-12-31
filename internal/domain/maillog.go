package domain

type MailLog struct {
	Uuid   string `bson:"uuid" json:"uuid"`
	Status string `bson:"status" json:"status"`
}
