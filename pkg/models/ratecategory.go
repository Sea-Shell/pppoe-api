package models

// RateCategory is a struct that represents a rating category
type RateCategory struct {
	CategoryID      *int64 `json:"category_id" db:"categoryId"`
	CategoryEventID int64  `json:"category_event_id" db:"categoryEventId"`
	CategoryName    string `json:"category_name" db:"categoryName"`
}
