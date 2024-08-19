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
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    router := gin.Default()

    router.POST("/login", func(c *gin.Context) {
        var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hashedPassword string
	var id int
	row := db.QueryRow("SELECT id, password FROM users WHERE username = ?", user.Username, user.Password)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	time64 := time.Now().Unix()
	timeInt := strconv.FormatInt(time64, 10)
	token := user.Username + user.Password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])

	c.SetCookie("auth_token", hashedToken, 3600, "/", "chatui.com", false, true)

	c.JSON(http.StatusOK, gin.H{"id": id, "token": hashedToken})
    })

    router.POST("/register", func(c *gin.Context) {
        var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO users (username, password, phone) VALUES (?, ?, ?)", user.Username, hashedPassword, user.Phone)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}