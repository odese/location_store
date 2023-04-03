package postgre

import (
	"context"
	"location_store/pkg/infrastructure/config"
	log "location_store/pkg/infrastructure/logger"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Init, initializes postgre connections.
func Init() {
	connectToPostgreSQL()
}

// Connection, represents a conn to PostgreSQL database.
var conn *pgxpool.Pool

// connectToPostgreSQL, connects to PostgreSQL database and initiates connection instance.
func connectToPostgreSQL() {
	var err error

	server := config.Call().GetString("postgre.server")
	port := config.Call().GetString("postgre.port")
	username := config.Call().GetString("postgre.username")
	password := config.Call().GetString("postgre.password")
	dbName := config.Call().GetString("postgre.dbname")

	connectionStr := "postgres://" + username + ":" + password + "@" + server + ":" + port + "/" + dbName + "?sslmode=disable"

	conn, err = pgxpool.New(context.Background(), connectionStr)
	if err != nil {
		log.Fatal().Err(err).Str("Connection Str", connectionStr).Msg("Error on connecting to PostgreSQL.")
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal().Err(err).Str("Connection Str", connectionStr).Msg("Error on pinging PostgreSQL.")
	} else {
		log.Info().Msg("PostgreSQL connection established.")
	}
}
