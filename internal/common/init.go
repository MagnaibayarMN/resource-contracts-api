package common

import (
	"log"

	"github.com/joho/godotenv"
)

func InitConfigurations(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
