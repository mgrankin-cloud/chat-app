package logic

type Notification struct {
    UserID  string `json:"user_id"`
    Message string `json:"message"`
}

func SendNotification(notification Notification) error {
    // написать логику отправки уведомления ( через Firebase Cloud Messaging (FCM) или другой сервис )
    return nil
}