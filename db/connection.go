package db

import (
	"database/sql"
	"fmt"
	"log"
)

type ConfigDB struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	MaxConns     int
	MaxIdleConns int
}

func OpenDbConnection(cfg ConfigDB) (*sql.DB, error) {
	log.Println("=> open db connection")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
