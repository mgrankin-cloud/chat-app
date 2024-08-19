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
        var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	if message.ImageBase64 != "" {
		imgData, err := base64.StdEncoding.DecodeString(message.ImageBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid base64 image"})
			return
		}

		imgPath := "./uploads/" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
		err = ioutil.WriteFile(imgPath, imgData, 0644)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}

		relPath, err := filepath.Rel("./uploads", imgPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get relative path"})
			return
		}

		message.ImageURL = "/" + relPath
	}

	result, err := db.Exec("INSERT INTO messages (channel_id, user_id, user_name, message, created_at, image_url) VALUES (?, ?, ?, ?, ?, ?)", message.ChannelID, message.UserID, message.UserName, message.Text, message.CreatedAt, message.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "photo_url": message.ImageURL})
    })

    router.GET("/chat/history", func(c *gin.Context) {
        channelID, err := strconv.Atoi(c.Query("channelID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channelID"})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 100
	}

	lastMessageID, err := strconv.Atoi(c.Query("lastMessageID"))
	if err != nil || lastMessageID < 0 {
		lastMessageID = 0
	}

	rows, err := db.Query("SELECT m.id, channel_id, user_id, u.username AS user_name, message, created_at, image_url FROM messages m LEFT JOIN users u ON u.id = m.user_id WHERE channel_id = ? AND m.id > ? ORDER BY m.id ASC LIMIT ?", channelID, lastMessageID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.ID, &message.ChannelID, &message.UserID, &message.UserName, &message.Text, &message.CreatedAt, &message.ImageURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}