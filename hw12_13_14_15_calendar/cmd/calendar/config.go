package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Storage StorageConf
	Logger  LoggerConf
	Service ServiceConf
}

func (c *Config) Validate() {
	if c.Service.Host == "" {
		log.Fatal("Host is not defined in the config file")
	}
	if len(c.Service.Port) == 0 { // todo: проверять, что содержит только цифры
		log.Fatal("Port is not defined in the config file")
	}
	if c.Storage.Type != "inmemory" && c.Storage.Type != "postgres" { // todo: вынести типы сторов в константу или enum
		log.Fatal("Incorrect storage type")
	}
}

type LoggerConf struct {
	Level string
	Path  string
}

type ServiceConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type     string
	Postgres StoragePostgresConf
}

type StoragePostgresConf struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewConfig(path string) (conf Config) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("error occurred during the reading the config file: %w", err))
	}

	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(fmt.Errorf("unable to decode into config struct,: %w", err))
	}

	return
}
