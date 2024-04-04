package models

// EventItem is a struct that represents an event with a rating associated with it
type EventItem struct {
	EventID          *int64 `json:"event_id" db:"eventId"`
	EventName        string `json:"event_name" db:"eventName"`
	EventCreated     string `json:"event_created" db:"eventCreated"`
	EventStart       string `json:"event_start" db:"eventStart"`
	EventEnd         string `json:"event_end" db:"eventEnd"`
	EventLocation    string `json:"event_location" db:"eventLocation"`
	EventDescription string `json:"event_description" db:"eventDescription"`
	EventOrganizer   string `json:"event_organizer" db:"eventOrganizer"`
}

// EventItemNoID is a struct that represents an event without id
type EventItemNoID struct {
	EventName        string `json:"event_name" db:"eventName"`
	EventCreated     string `json:"event_created" db:"eventCreated"`
	EventStart       string `json:"event_start" db:"eventStart"`
	EventEnd         string `json:"event_end" db:"eventEnd"`
	EventLocation    string `json:"event_location" db:"eventLocation"`
	EventDescription string `json:"event_description" db:"eventDescription"`
	EventOrganizer   string `json:"event_organizer" db:"eventOrganizer"`
}
