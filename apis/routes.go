package apis

import (
	"database/sql"

	"github.com/gorilla/mux"
)

//DBInstance stores db connection
type DBInstance struct {
	db *sql.DB
}

//RegisterRoutes defines all routes used in the app
func RegisterRoutes(r *mux.Router, db *sql.DB) {

	dbInstance := &DBInstance{db: db}

	r.HandleFunc("/email-senders", dbInstance.addEmailSender).Methods("POST")
	r.HandleFunc("/email-senders", dbInstance.getEmailSenders).Methods("GET")
	r.HandleFunc("/email-senders", dbInstance.deleteEmailSender).Methods("DELETE")
}
