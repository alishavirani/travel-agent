package models

type Config struct {
	ServerPort string `json:"serverPort"`
	SQL        SQL    `json:"sql"`
}

type SQL struct {
	DBType   string `json:"dbType"`
	Username string `json:"username"`
	Password string `json:"passoword"`
	DB       string `json:"db"`
}
