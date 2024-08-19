package models

import "time"

type Message struct {
    ID         string
    SenderID   string
    ReceiverID string
    Content    string
    Timestamp  time.Time
}