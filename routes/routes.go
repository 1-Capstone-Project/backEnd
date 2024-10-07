package routes

import (
	"database/sql"
	"gitmate/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	// CORS 설정
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 라우트 설정
	router.GET("/company_info", controllers.GetCompanyInfo(db))
	router.POST("/schedules", controllers.AddSchedule(db))
	router.GET("/schedules", controllers.GetSchedules(db))
	router.GET("/posts", controllers.GetPosts(db)) // GET 핸들러 추가

	return router
}
