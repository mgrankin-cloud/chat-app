cleanup() {
    echo "Killing procs..."
    kill $SERVER_PID
    kill $FLUTTER_PID
    exit
}

trap cleanup SIGINT SIGHUP

# Запуск микросервисов
./s.sh &
SERVER_PID=$!

# Запуск Flutter-приложения
cd /app/lib
flutter run -d <device_id> --web-port=3000 --dart-define=BASE_URL=http://localhost:8080 &
FLUTTER_PID=$!

wait $SERVER_PID
wait $FLUTTER_PID