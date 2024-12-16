package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"

	ssov5 "github.com/mgrankin-cloud/messenger/contract/gen/go/message"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()

	client := ssov5.NewMessageClient(grpcConn)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			processMessage(client, message)
		}
	}()

	for {
		select {
		case <-r.Context().Done():
			log.Println("connection closed")
			return
		}
	}
}

func processMessage(client ssov5.MessageClient, message []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &ssov5.CreateMessageRequest{
		Content: string(message),
		// мб другие поля добавить (received_by, created_by)
	}

	resp, err := client.CreateMessage(ctx, req)
	if err != nil {
		log.Printf("could not create message: %v", err)
		return
	}

	log.Printf("Message created with ID: %d", resp.MessageId)
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
