package config

import "os"

type Config struct {
	Port          string
	MySQLDSN      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		MySQLDSN:      getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/task_management?parseTime=true"),
		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
