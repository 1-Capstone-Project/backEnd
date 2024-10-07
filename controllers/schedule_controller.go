package controllers

import (
	"database/sql"
	"gitmate/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddSchedule(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSchedule models.Schedule

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
}

func GetSchedules(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := "SELECT id, title, description, schedule_date, COALESCE(start_time, ''), COALESCE(end_time, ''), COALESCE(img_url, '') FROM schedules"

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Database query error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var schedules []models.Schedule
		for rows.Next() {
			var schedule models.Schedule
			if err := rows.Scan(&schedule.ID, &schedule.Title, &schedule.Description, &schedule.ScheduleDate, &schedule.StartTime, &schedule.EndTime, &schedule.ImgURL); err != nil {
				log.Printf("Row scan error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			schedules = append(schedules, schedule)
		}

		c.JSON(http.StatusOK, schedules)
	}
}
