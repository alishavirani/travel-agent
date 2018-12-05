package apis

import (
	"net/http"
)

func welcomeRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO!!!"))
}
