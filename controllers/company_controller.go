package controllers

import (
	"database/sql"
	"fmt"
	"gitmate/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCompanyInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var companyInfos []models.CompanyInfo
		for rows.Next() {
			var companyInfo models.CompanyInfo
			if err := rows.Scan(&companyInfo.ID, &companyInfo.CompanyName, &companyInfo.HeadquartersLocation, &companyInfo.Industry, &companyInfo.Welfare, &companyInfo.RecruitmentMethod, &companyInfo.Requirements, &companyInfo.ImageURL); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			companyInfos = append(companyInfos, companyInfo)
		}

		c.JSON(http.StatusOK, companyInfos)
	}
}
