package models

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
