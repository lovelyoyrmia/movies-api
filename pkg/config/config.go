package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GatewayAddress string `mapstructure:"GATEWAY_ADDRESS"`
	DBUrl          string `mapstructure:"DB_URL"`
}

var ConfGlob Config

func LoadConfig() (config Config, err error) {
	godotenv.Load(".env")
	path := "."

	if os.Getenv("ENV") == "dev" {

		viper.AddConfigPath(path)
		viper.SetConfigName("dev")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		err = viper.Unmarshal(&config)

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Info().Msg("Load Development Environment...")
		return
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	log.Info().Msg("Load Production Environment...")
	return
}
