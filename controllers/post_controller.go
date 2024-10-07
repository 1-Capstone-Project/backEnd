package controllers

import (
	"database/sql"
	"gitmate/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPosts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, title, description, img_url FROM posts")
		if err != nil {
			log.Printf("Database query error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.ImgURL); err != nil {
				log.Printf("Row scan error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			posts = append(posts, post)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Rows iteration error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, posts)
	}
}
