package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type LinkServiceConfig struct {
	Address         string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

type LinkServiceConfigForStandaloneDB struct {
	Address     string
	BaseURL     string
	DatabaseDSN string
}

type LinkServiceConfigWithFileStorage struct {
	Address         string
	BaseURL         string
	FileStoragePath string
}

func (lsc *LinkServiceConfig) Parse() {
	// Считываем конфигурацию с помощью флагов
	flag.StringVar(&lsc.Address, "a", "localhost:8080", "server address")
	flag.StringVar(&lsc.BaseURL, "b", "http://localhost:8080", "base url")
	flag.StringVar(&lsc.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&lsc.DatabaseDSN, "d", "", "database dsn")
	flag.Parse()

	// Считываем конфигурацию с помощью переменных окружения
	var envCfg LinkServiceConfig
	err := env.Parse(&envCfg)
	if err != nil {
		log.Printf("Error parsing linservice config: %s\n", err)
		return
	}

	// Если в переменных окружения переданы какие-то значения, перезапишем конфигурацию с их помощью
	if envCfg.FileStoragePath != "" {
		lsc.FileStoragePath = envCfg.FileStoragePath
	}
	if envCfg.Address != "" {
		lsc.Address = envCfg.Address
	}
	if envCfg.BaseURL != "" {
		lsc.BaseURL = envCfg.BaseURL
	}
	if envCfg.DatabaseDSN != "" {
		lsc.DatabaseDSN = envCfg.DatabaseDSN
	}
}

func (lscfsadb *LinkServiceConfigForStandaloneDB) GetConfig(config LinkServiceConfig) {
	lscfsadb.Address = config.Address
	lscfsadb.BaseURL = config.BaseURL
	lscfsadb.DatabaseDSN = config.DatabaseDSN
}

func (lscwdb *LinkServiceConfigWithFileStorage) GetConfig(config LinkServiceConfig) {
	lscwdb.Address = config.Address
	lscwdb.BaseURL = config.BaseURL
	lscwdb.FileStoragePath = config.FileStoragePath
}
