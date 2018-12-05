package db

import (
	"database/sql"
	"log"
	"travel-agent-backend/models"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToMySql(config models.Config) *sql.DB {
	connectionString := config.SQL.Username + ":" + config.SQL.Password + "@/" + config.SQL.DB
	db, err := sql.Open(config.SQL.DBType, connectionString)
	if err != nil {
		log.Panic("Err in conecting to MySQL", err)
	}
	return db
}
