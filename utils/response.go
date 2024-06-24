package utils

import (
	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, statusCode int, payload interface{}) {
	c.JSON(statusCode, payload)
}

func ErrorJSON(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
