package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUser string
	DatabasePsw  string
	DatabaseURL  string
	DatabasePORT string
	DatabaseName string
	Host         string
	HostPort     string
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		DatabaseUser: os.Getenv("DBUSER"),
		DatabasePsw:  os.Getenv("DBPSW"),
		DatabaseURL:  os.Getenv("DBURL"),
		DatabasePORT: os.Getenv("DBPORT"),
		DatabaseName: os.Getenv("DBNAME"),
		Host:         os.Getenv("HOST"),
		HostPort:     os.Getenv("HOST_PORT"),
	}, nil

}
