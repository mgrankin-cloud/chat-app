package models

import "time"

type Message struct {
	ID         int64
	Content    string
	CreatedAt  time.Time
	CreatedBy  int64
	ReplyToID  int64
	ReceivedBy int64
	ReceivedAt time.Time
}
