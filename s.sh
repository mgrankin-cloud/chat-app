cleanup() {
    echo "Killing procs..."
    kill $AUTH_PID
    kill $MESSAGE_PID
    kill $CHAT_PID
    kill $MEDIA_PID
    kill $NOTIFICATIONS_PID
    kill $USER_PID
    kill $MODELS_PID
    exit
}

trap cleanup SIGINT SIGHUP

# Запуск сервиса авторизации
go run ./Chat-App/messenger/sso-service/cmd/main.go &
AUTH_PID=$!

# Запуск сервиса сообщений
go run ./chat-app/messenger/message-service/cmd/main.go &
MESSAGE_PID=$!

# Запуск сервиса чатов
go run ./chat-app/messenger/chat-service/cmd/main.go &
CHAT_PID=$!

go run ./chat-app/messenger/media-service/cmd/main.go &
MEDIA_PID=$!

go run ./chat-app/messenger/notification-service/cmd/main.go &
NOTIFICATIONS_PID=$!

go run ./chat-app/messenger/user-service/cmd/main.go &
USER_PID=$!

go run ./chat-app/messenger/models-service/cmd/main.go &
MODELS_PID=$!

wait $AUTH_PID
wait $MESSAGE_PID
wait $CHAT_PID
wait $MEDIA_PID
wait $NOTIFICATIONS_PID
wait $USER_PID
wait $MODELS_PID