package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"encoding/json"

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
	rows, err := db.Query("SELECT id, company_name, headquarters_location, industry, welfare, recruitment_method, requirements FROM company_info")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var companyInfos []CompanyInfo
	for rows.Next() {
		var companyInfo CompanyInfo
		if err := rows.Scan(&companyInfo.ID, &companyInfo.CompanyName, &companyInfo.HeadquartersLocation, &companyInfo.Industry, &companyInfo.Welfare, &companyInfo.RecruitmentMethod, &companyInfo.Requirements); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		companyInfos = append(companyInfos, companyInfo)
	}

	// JSON 정렬 및 변환
	companyInfosJSON, err := json.MarshalIndent(companyInfos, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", companyInfosJSON)
}
