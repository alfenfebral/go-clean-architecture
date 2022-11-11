package config

import (
	"github.com/joho/godotenv"
)

// LoadConfig - load environment config
func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
