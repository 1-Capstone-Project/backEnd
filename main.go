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

// CompanyInfo struct
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

	// 서버 실행
	router.Run(":8080")
}

// getCompanyInfo 핸들러 함수
func getCompanyInfo(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

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
