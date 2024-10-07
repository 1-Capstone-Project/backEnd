package models

type Schedule struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ScheduleDate string `json:"schedule_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	ImgURL       string `json:"img_url"`
}
