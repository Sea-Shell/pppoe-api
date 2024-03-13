package models

// common models used in the application

// ResponsePayload represents the response payload
type ResponsePayload struct {
	TotalItemCount int         `json:"total_item_count"`
	CurrentPage    int         `json:"current_page"`
	ItemLimit      int         `json:"item_limit"`
	TotalPages     int         `json:"total_pages"`
	Items          interface{} `json:"items"`
	NextPage       *string     `json:"next_page"`
	PrevPage       *string     `json:"prev_page"`
}

// Config represents the configuration file
type Config struct {
	Database Database `yaml:"database" json:"database"`
	General  General  `yaml:"general" json:"general"`
}

// Database represents the database configuration
type Database struct {
	File       string `yaml:"file" json:"file"`
	Connection string `yaml:"connection" json:"connection,omitempty"`
	Username   string `yaml:"username" json:"username,omitempty"`
	Password   string `yaml:"password" json:"password,omitempty"`
}

// General represents the general configuration
type General struct {
	Hostname   string   `yaml:"hostname" json:"hostname"`
	Schemes    []string `yaml:"schemes" json:"schemes"`
	ListenPort string   `yaml:"listen-port" json:"listen_port"`
	LogLevel   string   `yaml:"log-level" json:"log-level"`
}
