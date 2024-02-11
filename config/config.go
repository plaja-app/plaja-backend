package config

import "gorm.io/gorm"

// AppConfig holds the application config.
type AppConfig struct {
	DB  *gorm.DB
	Env *EnvVariables
}

// EnvVariables holds environment variables used in the application.
type EnvVariables struct {
	PostgresHost   string
	PostgresUser   string
	PostgresPass   string
	PostgresDBName string
	JWTSecret      string
}
