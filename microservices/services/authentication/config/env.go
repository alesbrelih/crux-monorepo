package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	Debug         bool          `mapstructure:"DEBUG"`
	AppPort       string        `mapstructure:"APP_PORT"`
	DatabaseUrl   string        `mapstructure:"DATABASE_URL"`
	JwtSecret     string        `mapstructure:"JWT_SECRET"`
	JwtAccessExp  time.Duration `mapstructure:"JWT_ACCESS_EXPIRATION"`
	JwtRefreshExt time.Duration `mapstructure:"JWT_REFRESH_EXPIRATION"`
	LogLevel      string        `mapstructure:"LOG_LEVEL"`
}

func GetEnvConfig(path, file string) *Env {
	// get enviroment variables
	var envConf Env
	viper.SetConfigFile(file)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read enviroment config. Err: %v", err)
	}
	if err := viper.Unmarshal(&envConf); err != nil {
		log.Fatalf("Failed to unmarshal enviroment variables. Err: %v", err)
	}
	return &envConf
}
