package server

import (
	"location_store/pkg/controllers"
	"location_store/pkg/infrastructure/config"
	log "location_store/pkg/infrastructure/logger"
	"net/http"
)

func Run() {
	setupHandlers()

	port := config.Call().GetString("host.port")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
	}
}

func setupHandlers() {
	http.HandleFunc("/csv", controllers.UploadCSV)
	http.HandleFunc("/listPlaces", controllers.ListPlaces)
}
