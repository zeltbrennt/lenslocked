package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/zeltbrennt/lenslocked/models"
)

type config struct {
	PSQL   models.PostgresConfig
	SMTP   models.SMTPConfig
	Server struct {
		Host   string
		Port   int
		Domain string
	}
}

func LoadFromEnv() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL.Host = os.Getenv("PSQL_HOST")
	cfg.PSQL.Port, err = strconv.Atoi(os.Getenv("PSQL_PORT"))
	if err != nil {
		log.Println("psql port")
		return cfg, err
	}
	cfg.PSQL.User = os.Getenv("PSQL_USERNAME")
	cfg.PSQL.Password = os.Getenv("PSQL_PASSWORD")
	cfg.PSQL.Dbname = os.Getenv("PSQL_DBNAME")
	cfg.PSQL.Sslmode = os.Getenv("PSQL_SSLMODE")

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println("smtp port")
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	cfg.Server.Host = os.Getenv("HOST")
	cfg.Server.Domain = os.Getenv("DOMAIN")
	cfg.Server.Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("server port")
		return cfg, err
	}
	return cfg, nil
}
