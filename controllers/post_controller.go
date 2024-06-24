package controllers

import (
	"database/sql"
	"gitmate/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newPost models.Post

		if err := c.ShouldBindJSON(&newPost); err != nil {
			log.Printf("JSON binding error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `INSERT INTO posts (title, description, img_url) VALUES (?, ?, ?)`
		result, err := db.Exec(query, newPost.Title, newPost.Description, newPost.ImgURL)
		if err != nil {
			log.Printf("Database insert error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			log.Printf("Getting last insert ID error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Post added successfully with ID: %d", id)
		c.JSON(http.StatusOK, gin.H{"message": "Post added successfully", "id": id})
	}
}
