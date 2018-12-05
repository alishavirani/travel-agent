package apis

import (
	"github.com/gorilla/mux"
)

//RegisterRoutes defines all routes used in the app
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", welcomeRoute).Methods("GET")
}
