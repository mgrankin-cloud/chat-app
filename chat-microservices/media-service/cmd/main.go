package main

import (
    "log"
    "net/http"
    "github.com/yourusername/media-service/internal/config"
    "github.com/yourusername/media-service/internal/logic"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    router := gin.Default()

    router.POST("/upload", func(c *gin.Context) {

        file, err := c.FormFile("image")
	    if err != nil {
		    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		    return
	    }

     	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		  err := os.Mkdir("./uploads", os.ModePerm)
		    if err != nil {
			    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			    return
		    }
	    }

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resizedImg := imaging.Fit(img, 64, 64, imaging.Lanczos)

	fileName := fmt.Sprintf("./uploads/%s.png", uuid.New().String())

	err = imaging.Save(resizedImg, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := fmt.Sprintf("/uploads/%s", filepath.Base(fileName))
	c.JSON(http.StatusOK, gin.H{"url": url})

        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := c.SaveUploadedFile(file, "uploads/"+file.Filename); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
    })

    router.GET("/media/:id", func(c *gin.Context) {
        mediaID := c.Param("id")
        // написать логику получения медиафайла по ID
        mediaFile, err := logic.GetMediaFileByID(mediaID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Media file not found"})
            return
        }

        c.File("uploads/" + mediaFile.Filename)
    })

    log.Printf("Starting server on port %d", cfg.Port)
    if err := router.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}