package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CompanyInfo struct {
	ID                   int    `json:"id"`
	CompanyName          string `json:"company_name"`
	HeadquartersLocation string `json:"headquarters_location"`
	Industry             string `json:"industry"`
	Welfare              string `json:"welfare"`
	RecruitmentMethod    string `json:"recruitment_method"`
	Requirements         string `json:"requirements"`
	ImageURL             string `json:"image_url"`
}

type Schedule struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ScheduleDate string `json:"schedule_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ImgURL       string `json:"img_url"`
}

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImgURL      string `json:"img_url"`
}

var db *sql.DB

func main() {
	var err error

	dsn := "admin:gitmate1234@tcp(gitmate-database.cbimo8eqih5y.ap-northeast-2.rds.amazonaws.com:3306)/gitmate_db?charset=utf8"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Database connection established")

	router := gin.Default()

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

	router.GET("/company_info", getCompanyInfo)
	router.POST("/schedules", addSchedule)
	router.GET("/schedules", getSchedules)
	router.POST("/posts", addPost)

	router.Run(":8080")
}

func getCompanyInfo(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT id, company_name, headquarters_location, industry, welfare, recruitment_method, requirements, image_url FROM company_info LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var companyInfos []CompanyInfo
	for rows.Next() {
		var companyInfo CompanyInfo
		if err := rows.Scan(&companyInfo.ID, &companyInfo.CompanyName, &companyInfo.HeadquartersLocation, &companyInfo.Industry, &companyInfo.Welfare, &companyInfo.RecruitmentMethod, &companyInfo.Requirements, &companyInfo.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		companyInfos = append(companyInfos, companyInfo)
	}

	c.JSON(http.StatusOK, companyInfos)
}

func addSchedule(c *gin.Context) {
	var newSchedule Schedule

	if err := c.ShouldBindJSON(&newSchedule); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO schedules (title, description, schedule_date, start_time, end_time, img_url) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, newSchedule.Title, newSchedule.Description, newSchedule.ScheduleDate, newSchedule.StartTime, newSchedule.EndTime, newSchedule.ImgURL)
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

	log.Printf("Schedule added successfully with ID: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Schedule added successfully", "id": id})
}

func getSchedules(c *gin.Context) {
	query := "SELECT id, title, description, schedule_date, COALESCE(start_time, ''), COALESCE(end_time, ''), COALESCE(img_url, '') FROM schedules"

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ID, &schedule.Title, &schedule.Description, &schedule.ScheduleDate, &schedule.StartTime, &schedule.EndTime, &schedule.ImgURL); err != nil {
			log.Printf("Row scan error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

func addPost(c *gin.Context) {
	var newPost Post

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
