package main

import (
	"location_store/pkg/infrastructure/config"
	"location_store/pkg/infrastructure/db/postgre"
	log "location_store/pkg/infrastructure/logger"
	"location_store/pkg/server"
)

func init() {
	log.SetupLogger()
	config.Init()
	postgre.Init()
}

func main() {
	server.Run()
}
