package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
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

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
