// internal/config/config.go
package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	MySQL struct {
		DNS string
	}
	Redis struct {
		Addr         string
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
	}
	Key struct {
		Salt string
	}
	Timeout struct {
		DelayHeartbeat   int
		HeartbeatHz      int
		HeartbeatMaxTime int
		RedisOnlineTime  int
	}
	Port struct {
		Server string
		UDP    int
	}
}

var Config *AppConfig

func LoadConfig(path string) {
	viper.SetConfigName("app")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	Config = &AppConfig{}

	if err := viper.Unmarshal(Config); err != nil {
		log.Fatal("Error unmarshaling config: ", err)
	}

	log.Println("Configuration loaded successfully.")
}
