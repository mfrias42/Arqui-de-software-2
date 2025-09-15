package config

import "time"

const (
	MySQLHost     = "mysql"
	MySQLPort     = "3306"
	MySQLDatabase = "users"
	MySQLUsername = "root"
	MySQLPassword = "root"

	CacheDuration = 30 * time.Second

	MemcachedHost = ""
	MemcachedPort = "11211"

	JWTKey      = "ThisIsAnExampleJWTKey!"
	JWTDuration = 24 * time.Hour
)
