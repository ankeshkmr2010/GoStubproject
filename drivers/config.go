package drivers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Configuration struct {
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

func LoadConfiguration() *Configuration {
	config := &Configuration{}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if len(os.Getenv("DB_HOST")) != 0 {
		config.DbHost = os.Getenv("DB_HOST")
	} else {
		log.Fatal("DB_HOST is not set")
	}
	if len(os.Getenv("DB_PORT")) != 0 {
		config.DbPort = os.Getenv("DB_PORT")
	} else {
		log.Fatal("DB_PORT is not set")
	}

	if len(os.Getenv("DB_USER")) != 0 {
		config.DbUser = os.Getenv("DB_USER")
	} else {
		log.Fatal("DB_USER is not set")
	}
	if len(os.Getenv("DB_PASS")) != 0 {
		config.DbPass = os.Getenv("DB_PASS")
	} else {
		log.Fatal("DB_PASS is not set")
	}

	if len(os.Getenv("DB_DATABASE")) != 0 {
		config.DbName = os.Getenv("DB_DATABASE")
	} else {
		log.Fatal("DB_DATABASE is not set")
	}

	return config
}
