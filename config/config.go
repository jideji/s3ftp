package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config is a config.
type Config struct {
	Ftp Ftp
	S3  S3
}

type Ftp struct {
	Host     string
	Port     int
	Username string
	Password string
}

type S3 struct {
	BucketName string
}

func getEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Errorf("No environment variable %s defined", key))
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		nValue, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return nValue
	}
	return defaultValue
}

func LoadConfig() *Config {
	return &Config{
		Ftp: Ftp{
			Host:     getEnv("FTP_HOST", "localhost"),
			Port:     getEnvInt("FTP_PORT", 21),
			Username: mustGetEnv("FTP_USER"),
			Password: mustGetEnv("FTP_PASS"),
		},
		S3: S3{
			BucketName: mustGetEnv("S3_BUCKET_NAME"),
		},
	}
}
