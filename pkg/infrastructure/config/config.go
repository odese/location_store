package config

import (
	"strings"

	log "location_store/pkg/infrastructure/logger"
	"location_store/pkg/utils"
	"github.com/spf13/viper"
)

// Represents the instance of config file.
var config *viper.Viper

// Init, reads yml file and initializes config instance.
func Init() {
	config = viper.New()
	config.SetConfigType("yml")
	config.AddConfigPath("configs")
	config.SetConfigName(getNameOfConfigFile())

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error on reading config file.")
	}
}

// getNameOfConfigFile, chooses proper config file name according to work environment.
func getNameOfConfigFile() (fileName string) {
	log.Info().Str("WORK_ENV", utils.WorkEnvironment).Send()
	fileName = strings.ToLower(utils.WorkEnvironment)
	return fileName
}

// Call, returns instance of config file.
func Call() *viper.Viper {
	return config
}

// InitForTest, reads yml file and initializes config instance for testing.
func InitForTest() {
	config = viper.New()
	config.SetConfigType("yml")
	config.AddConfigPath("c:/workspace/go/src/location_store/configs")
	config.SetConfigName(getNameOfConfigFile())

	err := config.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error on reading config file.")
	}
}