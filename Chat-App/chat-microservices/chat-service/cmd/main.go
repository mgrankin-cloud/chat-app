package main

import (
    "log"
    "net/http"
    "github.com/mgrankin-cloud/chat-service/internal/config"
    "github.com/mgrankin-cloud/chat-service/internal/logic"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    router := gin.Default()

    router.POST("/chat/send", func(c *gin.Context) {
        // написать логику отправки сообщения
        c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
    })

    router.GET("/chat/history", func(c *gin.Context) {
        // написать логику получения истории сообщений
        c.JSON(http.StatusOK, gin.H{"history": "Chat history"})
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}