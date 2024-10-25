package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		PORT                                 string `env:"PORT"                                          envDefault:"3000"`
		APP_NAME                             string `env:"APP_NAME"                                      envDefault:"Boilerplate"`
		APP_URL                              string `env:"APP_URL"`
		GITHUB_API_TOKEN                     string `env:"GITHUB_API_TOKEN,required"`
		GITHUB_API_BASE_URL                  string `env:"GITHUB_API_BASE_URL,required"`
		MONGO_DB_CONFIG                      MONGO_DB_CONFIG
		RABBIT_MQ_CONFIG                     RABBIT_MQ_CONFIG
		REPOSITORIES_PATH                    string `env:"REPOSITORIES_PATH,required"`
		BASE_DOMAIN                          string `env:"BASE_DOMAIN,required"`
		TRAEFIK_NETWORK_NAME                 string `env:"TRAEFIK_NETWORK_NAME,required"`
		MAX_DEPLOYING_STATUS_TIME_IN_MINUTES int    `env:"MAX_DEPLOYING_STATUS_TIME_IN_MINUTES,required"`
	}
	DB_CONFIG struct {
		DB_NAME       string `env:"DB_NAME"`
		HOST          string `env:"DB_HOST"`
		PORT          string `env:"DB_PORT"`
		USER          string `env:"DB_USER"`
		PASSWORD      string `env:"DB_PASSWORD"`
		DEBUG_ENABLED string `env:"DEBUG_ENABLED"`
	}
	MONGO_DB_CONFIG struct {
		URL     string `env:"MONGO_DB_URL,required"`
		DB_NAME string `env:"MONGO_DB_NAME,required"`
	}
	RABBIT_MQ_CONFIG struct {
		URL                                      string `env:"RABBIT_MQ_CONNECTION_URL,required"`
		EXCHANGE                                 string `env:"RABBIT_MQ_EXCHANGE,required"`
		REPOSITORY_PULL_ROUTING_KEY              string `env:"RABBIT_MQ_REPOSITORY_PULL_ROUTING_KEY,required"`
		REPOSITORY_PULL_QUEUE                    string `env:"RABBIT_MQ_REPOSITORY_PULL_QUEUE,required"`
		REPOSITORY_BUILD_ROUTING_KEY             string `env:"RABBIT_MQ_REPOSITORY_BUILD_ROUTING_KEY,required"`
		REPOSITORY_BUILD_QUEUE                   string `env:"RABBIT_MQ_REPOSITORY_BUILD_QUEUE,required"`
		REPOSITORY_RUN_ROUTING_KEY               string `env:"RABBIT_MQ_REPOSITORY_RUN_ROUTING_KEY,required"`
		REPOSITORY_RUN_QUEUE                     string `env:"RABBIT_MQ_REPOSITORY_RUN_QUEUE,required"`
		REPOSITORY_STOP_ROUTING_KEY              string `env:"RABBIT_MQ_REPOSITORY_STOP_ROUTING_KEY,required"`
		REPOSITORY_STOP_QUEUE                    string `env:"RABBIT_MQ_REPOSITORY_STOP_QUEUE,required"`
		RABBIT_MQ_REPOSITORY_PRE_RUN_ROUTING_KEY string `env:"RABBIT_MQ_REPOSITORY_PRE_RUN_ROUTING_KEY,required"`
		RABBIT_MQ_REPOSITORY_PRE_RUN_QUEUE       string `env:"RABBIT_MQ_REPOSITORY_PRE_RUN_QUEUE,required"`
	}
)

var AppConfig Config

func Init() error {
	err := parseConfigFile()
	if err != nil {
		return err
	}
	fmt.Println("Configuration loaded")
	return nil
}

func parseConfigFile() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file ", err.Error())
	}

	return env.Parse(&AppConfig)
}
