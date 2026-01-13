package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
	sslmode  string
}

func (pc PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", pc.host, pc.port, pc.user, pc.password, pc.dbname, pc.sslmode)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		"localhost",
		5432,
		"baloo",
		"junglebook",
		"lenslocked",
		"disable",
	}
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	return db, nil
}
