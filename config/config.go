/**
 * Copyright 2023 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env                    string `mapstructure:"ENV_NAME"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
	ApiKeyName             string `mapstructure:"COINBASE_CLOUD_API_KEY_NAME"`
	ApiPrivateKey          string `mapstructure:"COINBASE_CLOUD_API_PRIVATE_KEY"`
	AppUrl                 string `mapstructure:"APP_URL"`
	S3Enabled              bool   `mapstructure:"S3_ENABLED"`
	BucketName             string `mapstructure:"BUCKET_NAME"`
	PresignedUrlExpiration int64  `mapstructure:"PRESIGNED_URL_EXPIRATION"`
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
