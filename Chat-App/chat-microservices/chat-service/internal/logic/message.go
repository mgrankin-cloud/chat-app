package logic

import (
    "time"
    "github.com/mgrankin-cloud/auth-service/internal/models"
)

func SendMessage(senderID, receiverID string, content string) error {
    message := &models.Message{
        SenderID:   senderID,
        ReceiverID: receiverID,
        Content:    content,
        Timestamp:  time.Now(),
    }
    // тут сделать сохранение сообщения в базу данных
    return nil
}