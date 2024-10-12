package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type (
	Config struct {
		PORT                 string `env:"PORT"                 envDefault:"3000"`
		APP_NAME             string `env:"APP_NAME"             envDefault:"Boilerplate"`
		APP_URL              string `env:"APP_URL"`
		GITHUB_API_TOKEN     string `env:"GITHUB_API_TOKEN"`
		GITHUB_API_BASE_URL  string `env:"GITHUB_API_BASE_URL"`
		DB_DRIVER_NAME       string `env:"DB_DRIVER_NAME"`
		DB_CONFIG            DB_CONFIG
		MONGO_DB_CONFIG      MONGO_DB_CONFIG
		RABBIT_MQ_CONFIG     RABBIT_MQ_CONFIG
		REPOSITORIES_PATH    string `env:"REPOSITORIES_PATH"`
		BASE_DOMAIN          string `env:"BASE_DOMAIN"`
		TRAEFIK_NETWORK_NAME string `env:"TRAEFIK_NETWORK_NAME"`
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
		URL     string `env:"MONGO_DB_URL"`
		DB_NAME string `env:"MONGO_DB_NAME"`
	}
	RABBIT_MQ_CONFIG struct {
		URL                          string `env:"RABBIT_MQ_CONNECTION_URL"`
		EXCHANGE                     string `env:"RABBIT_MQ_EXCHANGE"`
		REPOSITORY_PULL_ROUTING_KEY  string `env:"RABBIT_MQ_REPOSITORY_PULL_ROUTING_KEY"`
		REPOSITORY_PULL_QUEUE        string `env:"RABBIT_MQ_REPOSITORY_PULL_QUEUE"`
		REPOSITORY_BUILD_ROUTING_KEY string `env:"RABBIT_MQ_REPOSITORY_BUILD_ROUTING_KEY"`
		REPOSITORY_BUILD_QUEUE       string `env:"RABBIT_MQ_REPOSITORY_BUILD_QUEUE"`
		REPOSITORY_RUN_ROUTING_KEY   string `env:"RABBIT_MQ_REPOSITORY_RUN_ROUTING_KEY"`
		REPOSITORY_RUN_QUEUE         string `env:"RABBIT_MQ_REPOSITORY_RUN_QUEUE"`
		REPOSITORY_STOP_ROUTING_KEY  string `env:"RABBIT_MQ_REPOSITORY_STOP_ROUTING_KEY"`
		REPOSITORY_STOP_QUEUE        string `env:"RABBIT_MQ_REPOSITORY_STOP_QUEUE"`
	}
)

var AppConfig Config

func Init(envFilePath string) error {
	err := parseConfigFile(envFilePath, "main")
	if err != nil {
		return err
	}
	fmt.Println(AppConfig)

	fmt.Println("Configuration loaded")
	return nil
}

func setFromEnv(config *Config) {
	config.PORT = viper.GetString("PORT")
	config.APP_NAME = viper.GetString("APP_NAME")
	config.APP_URL = viper.GetString("APP_URL")
	config.GITHUB_API_TOKEN = viper.GetString("GITHUB_API_TOKEN")
	config.GITHUB_API_BASE_URL = viper.GetString("GITHUB_API_BASE_URL")
	config.DB_CONFIG.DB_NAME = viper.GetString("DB_NAME")
	config.DB_CONFIG.HOST = viper.GetString("DB_HOST")
	config.DB_CONFIG.PORT = viper.GetString("DB_PORT")
	config.DB_CONFIG.PASSWORD = viper.GetString("DB_PASSWORD")
	config.DB_DRIVER_NAME = viper.GetString("DB_DRIVER_NAME")
	config.DB_CONFIG.USER = viper.GetString("DB_USER")
	config.DB_CONFIG.DEBUG_ENABLED = viper.GetString("DB_DEBUG_ENABLED")
	config.MONGO_DB_CONFIG.URL = viper.GetString("MONGO_DB_URL")
	config.MONGO_DB_CONFIG.DB_NAME = viper.GetString("MONGO_DB_NAME")
	config.RABBIT_MQ_CONFIG.EXCHANGE = viper.GetString("RABBIT_MQ_EXCHANGE")
	config.RABBIT_MQ_CONFIG.URL = viper.GetString("RABBIT_MQ_CONNECTION_URL")
	config.RABBIT_MQ_CONFIG.REPOSITORY_PULL_ROUTING_KEY = viper.GetString(
		"RABBIT_MQ_REPOSITORY_PULL_ROUTING_KEY",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE = viper.GetString(
		"RABBIT_MQ_REPOSITORY_PULL_QUEUE",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_ROUTING_KEY = viper.GetString(
		"RABBIT_MQ_REPOSITORY_BUILD_ROUTING_KEY",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_QUEUE = viper.GetString(
		"RABBIT_MQ_REPOSITORY_BUILD_QUEUE",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE = viper.GetString("RABBIT_MQ_REPOSITORY_RUN_QUEUE")
	config.RABBIT_MQ_CONFIG.REPOSITORY_RUN_ROUTING_KEY = viper.GetString(
		"RABBIT_MQ_REPOSITORY_RUN_ROUTING_KEY",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_STOP_QUEUE = viper.GetString(
		"RABBIT_MQ_REPOSITORY_STOP_QUEUE",
	)
	config.RABBIT_MQ_CONFIG.REPOSITORY_STOP_ROUTING_KEY = viper.GetString(
		"RABBIT_MQ_REPOSITORY_STOP_ROUTING_KEY",
	)
	config.REPOSITORIES_PATH = viper.GetString("REPOSITORIES_PATH")
	config.BASE_DOMAIN = viper.GetString("BASE_DOMAIN")
	config.TRAEFIK_NETWORK_NAME = viper.GetString("TRAEFIK_NETWORK_NAME")
}

func parseConfigFile(envFilePath, configName string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return err
	}
	return env.Parse(&AppConfig)
}
