package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var c config

type config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Username     string
	Password     string
	Name         string
	Host         string
	Port         uint64
	Encoding     string
	Maxconns     uint64
	Maxidleconns uint64
	Timeout      uint64
	Logmode      bool
}

// 初始化資料庫
func Init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.Unmarshal(&c)
}

func GetDatabaseConfig() DatabaseConfig {
	return c.Database
}
