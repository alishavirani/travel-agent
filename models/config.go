package models

type Config struct {
	ServerPort       string  `json:"serverPort"`
	SQL              SQL     `json:"sql"`
	Email            Email   `json:"email"`
	GmailAPIInterval float64 `json:"gmailApiInterval"`
}

type SQL struct {
	DBType   string `json:"dbType"`
	Username string `json:"username"`
	Password string `json:"passoword"`
	DB       string `json:"db"`
}

type Email struct {
	EmailID string `json:"emailId"`
	User    string `json:"user"`
}
