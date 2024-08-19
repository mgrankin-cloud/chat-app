package main

import (
    "log"
    "net/http"
    "github.com/mgrankin-cloud/notification-service/internal/config"
    "github.com/mgrankin-cloud/notification-service/internal/logic"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    router := gin.Default()

    router.POST("/notify", func(c *gin.Context) {
        var notification logic.Notification
        if err := c.ShouldBindJSON(&notification); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := logic.SendNotification(notification); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}