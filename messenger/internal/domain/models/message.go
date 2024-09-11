package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	ID         int64
	Content    string
	CreatedAt  *timestamppb.Timestamp
	CreatedBy  int64
	ReplyToID  int64
	ReceivedBy int64
	ReceivedAt *timestamppb.Timestamp
	Status string
	IsRead     bool
}
