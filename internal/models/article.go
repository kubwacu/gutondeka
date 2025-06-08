package models

import "time"

type Article struct {
	ContributorID int       `json:"contributor_id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Source        string    `json:"source"`
	Date          time.Time `json:"date"`
	Category      string    `json:"category"`
	FilePath      string    `json:"file_path,omitempty"`
}

type ArticleCount struct {
	MonthYear string `json:"month_year"`
	Count     int    `json:"count"`
}
