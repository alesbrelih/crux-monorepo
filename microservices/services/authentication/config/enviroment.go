package config

import "time"

type Enviroment struct {
	Debug         bool          `mapstructure:"DEBUG"`
	AppPort       string        `mapstructure:"APP_PORT"`
	DatabaseUrl   string        `mapstructure:"DATABASE_URL"`
	JwtSecret     string        `mapstructure:"JWT_SECRET"`
	JwtAccessExp  time.Duration `mapstructure:"JWT_ACCESS_EXPIRATION"`
	JwtRefreshExt time.Duration `mapstructure:"JWT_REFRESH_EXPIRATION"`
}
