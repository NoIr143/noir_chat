package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/noir143/noir_chat/src/configs"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg PostgresConfig) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
}

func NewPostgresStorage() (*sql.DB, error) {
	dbCfg := PostgresConfig{
		Host:     configs.EnvConfigs.DB_HOST,
		Port:     configs.EnvConfigs.DB_PORT,
		User:     configs.EnvConfigs.DB_USER,
		Password: configs.EnvConfigs.DB_PASSWORD,
		DBName:   configs.EnvConfigs.DB_NAME,
		SSLMode:  configs.EnvConfigs.DB_SSLMODE,
	}

	db, err := sql.Open("postgres", dbCfg.ConnString())
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
