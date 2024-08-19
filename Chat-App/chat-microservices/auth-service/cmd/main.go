package main

import (
	"fmt"
    "log"
    "net/http"
    "github.com/mgrankin-cloud/auth-service/internal/config"
    "github.com/mgrankin-cloud/auth-service/internal/logic"
    "github.com/gin-gonic/gin"
)

func main() {
    // Загрузка конфигурации
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Инициализация роутера
    router := gin.Default()

    router.POST("/login", func(c *gin.Context) {
        
        c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
    })

    router.POST("/register", func(c *gin.Context) {
        // Логика регистрации
        c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}