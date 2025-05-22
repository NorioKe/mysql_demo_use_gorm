package config

type DatabaseConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"dbname"`
	Charset   string `json:"charset"`
	ParseTime bool   `json:"parseTime"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
}
