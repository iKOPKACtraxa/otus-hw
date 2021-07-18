package main

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	HTTPServer HTTPServerConf
	Storage    StorageConf
}

type LoggerConf struct {
	File  string `toml:"logger.file"`
	Level string `toml:"logger.level"`
}

type HTTPServerConf struct {
	HostPort string `toml:"httpServer.hostPort"`
}

type StorageConf struct {
	Type    string `toml:"storage.type"`
	ConnStr string `toml:"storage.connStr"`
}

// NewConfig make a config from configFilePath.
func NewConfig() Config {
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	// черновик, удалю
	// fmt.Println(viper.ConfigFileUsed())
	// fmt.Println(viper.AllSettings())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return config
}
