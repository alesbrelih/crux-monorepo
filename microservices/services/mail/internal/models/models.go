package models

import "time"

const DEFAULT_MAIL_STATUS = "IN_QUEUE"

type Mail struct {
	Id        int64
	Reciever  string
	Subject   string
	Body      string
	Status    string
	CreatedAt time.Time
}
