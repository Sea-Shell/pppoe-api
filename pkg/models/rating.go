package models

// Rating is a struct that represents a rating event
type Rating struct {
	RatingID         *int64 `json:"rating_id" db:"ratingId"`
	RatingCategoryID int64  `json:"rating_category_id" db:"ratingCategoryId"`
	EventID          int64  `json:"event_id" db:"eventId"`
	RateName         string `json:"rate_name" db:"rateName"`
	RateScaleID      int64  `json:"rate_scale_id" db:"rateScaleId"`
	RateScaleValues  int64  `json:"rate_scale_values" db:"rateScaleValues"`
}
