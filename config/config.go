package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env           string `mapstructure:"ENV_NAME"`
	LogLevel      string `mapstructure:"LOG_LEVEL"`
	ApiKeyName    string `mapstructure:"COINBASE_CLOUD_API_KEY_NAME"`
	ApiPrivateKey string `mapstructure:"COINBASE_CLOUD_API_PRIVATE_KEY"`
	AppUrl        string `mapstructure:"APP_URL"`
}

func (a AppConfig) IsLocalEnv() bool {
	return a.Env == "local"
}

var appConfig *AppConfig

func Get() *AppConfig {
	return appConfig
}

func Setup(app *AppConfig) error {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("ENV_NAME", "local")
	viper.SetDefault("COINBASE_CLOUD_API_KEY_NAME", "NOT_SET")
	viper.SetDefault("COINBASE_CLOUD_API_PRIVATE_KEY", "NOT_SET")
	viper.SetDefault("APP_URL", "http://localhost")

	err := viper.ReadInConfig()
	if err != nil {
		log.Debugf("Missing env file %v", err)
	}

	err = viper.Unmarshal(&app)
	if err != nil {
		log.Debugf("cannot parse env file %v", err)
	}

	appConfig = app

	return nil
}
