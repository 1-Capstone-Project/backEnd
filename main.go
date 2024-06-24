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

// CompanyInfo struct 정의
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

// Schedule struct 정의
type Schedule struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ScheduleDate string `json:"schedule_date"` // `time.Time` 대신 string 사용
	StartTime    string `json:"start_time"`    // `time.Time` 대신 string 사용
	EndTime      string `json:"end_time"`      // `time.Time` 대신 string 사용
	ImgURL       string `json:"img_url"`
}

// Post struct 정의
type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImgURL      string `json:"img_url"`
}

var db *sql.DB

func main() {
	var err error

	// 데이터베이스 연결 정보 설정 (charset=utf8 추가)
	dsn := "admin:gitmate1234@tcp(gitmate-database.cbimo8eqih5y.ap-northeast-2.rds.amazonaws.com:3306)/gitmate_db?charset=utf8"

	// 데이터베이스 연결
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// 데이터베이스 핑 테스트
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	fmt.Println("Database connection established")

	// 라우터 설정
	router := gin.Default()

	// API 엔드포인트 설정
	router.GET("/company_info", getCompanyInfo)
	router.POST("/schedules", addSchedule)
	router.GET("/schedules", getSchedules) // GET /schedules 엔드포인트 추가
	router.POST("/posts", addPost)         // 새로운 엔드포인트 추가

	// 서버 실행
	router.Run(":8080")
}

// getCompanyInfo 핸들러 함수
func getCompanyInfo(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20") // 기본 limit 설정

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

	c.JSON(http.StatusOK, companyInfos) // 최종 결과 전송
}

// addSchedule 핸들러 함수
func addSchedule(c *gin.Context) {
	var newSchedule Schedule

	// 요청 바디에서 JSON 데이터를 파싱
	if err := c.ShouldBindJSON(&newSchedule); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// MySQL에 데이터 삽입
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

// getSchedules 핸들러 함수 추가
func getSchedules(c *gin.Context) {
	query := "SELECT id, title, description, schedule_date, start_time, end_time, IFNULL(img_url, '') FROM schedules"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var schedules []Schedule
	for rows.Next() {
		var schedule Schedule
		if err := rows.Scan(&schedule.ID, &schedule.Title, &schedule.Description, &schedule.ScheduleDate, &schedule.StartTime, &schedule.EndTime, &schedule.ImgURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules) // 최종 결과 전송
}

// addPost 핸들러 함수
func addPost(c *gin.Context) {
	var newPost Post

	// 요청 바디에서 JSON 데이터를 파싱
	if err := c.ShouldBindJSON(&newPost); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// MySQL에 데이터 삽입
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
