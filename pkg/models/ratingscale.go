package models

// RatingScale is a struct that represents a scale of a rating with rating type and pattern of rating
type RatingScale struct {
	ScaleID          *int64 `json:"scale_id" db:"scaleId"`
	ScaleName        string `json:"scale_name" db:"scaleName"`
	ScaleDescription string `json:"scale_description" db:"scaleDescription"`
	ScaleExample     string `json:"scale_example" db:"scaleExample"`
	ScalePattern     string `json:"scale_pattern" db:"scalePattern"`
}
