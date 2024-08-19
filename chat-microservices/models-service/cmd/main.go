package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/models-service/internal/config"
	"github.com/yourusername/models-service/internal/logic"
	"github.com/yourusername/models-service/internal/models"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := gin.Default()

	router.POST("/models/user", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := logic.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	})

	router.GET("/models/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		user, err := logic.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.PUT("/models/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.ID = userID
		if err := logic.UpdateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	router.DELETE("/models/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		if err := logic.DeleteUser(userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, nil)
	})

	log.Printf("Starting server on port %d", cfg.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
