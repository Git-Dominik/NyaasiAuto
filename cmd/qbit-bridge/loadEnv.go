package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Pass string
	User string
}

func loadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	eC := &Config{
		Pass: os.Getenv("qbPass"),
		User: os.Getenv("qbUsr"),
	}

	return eC, nil
}
