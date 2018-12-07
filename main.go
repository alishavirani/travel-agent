package main

import (
	"fmt"
	"net/http"
	"sync"
	"travel-agent-backend/apis"
	"travel-agent-backend/db"
	"travel-agent-backend/services"
	"travel-agent-backend/utils"

	"github.com/gorilla/mux"
)

func main() {
	//Load config
	path := "C:/Users/Alisha Virani/go/src/travel-agent-backend/config/config.development.json"
	config := utils.LoadConfig(path)

	//Connect to MySql
	db := db.ConnectToMySql(config)

	var wg sync.WaitGroup
	wg.Add(2)

	//Create router
	go func(w *sync.WaitGroup) {
		r := mux.NewRouter()
		apis.RegisterRoutes(r)
		port := ":" + config.ServerPort
		http.ListenAndServe(port, r)
		fmt.Println("Server is listening on port ", port)
	}(&wg)

	//Start service
	go func(w *sync.WaitGroup) {
		services.EmailReader(config, db)
	}(&wg)

	wg.Wait()
}
