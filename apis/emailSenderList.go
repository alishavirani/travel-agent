package apis

import (
	"net/http"
)

func (db *DBInstance) addEmailSender(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
}

func (db *DBInstance) getEmailSenders(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
}

func (db *DBInstance) deleteEmailSender(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
}
