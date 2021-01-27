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
	Storage     StorageConf
	Logger      LoggerConf
	Service     ServiceConf
	GRPCService ServiceConf
}

func (c *Config) Validate() []error {
	errors := make([]error, 0)
	if c.Service.Host == "" {
		errors = append(errors, fmt.Errorf("host is not defined in the config file"))
	}
	if len(c.Service.Port) == 0 {
		errors = append(errors, fmt.Errorf("port is not defined in the config file"))
	}
	if c.Storage.Type != "inmemory" && c.Storage.Type != "postgres" { // todo: вынести типы сторов в константу или enum
		errors = append(errors, fmt.Errorf("incorrect storage type"))
	}

	return errors
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
