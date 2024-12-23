package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type ConfigDB struct {
	Driver       string
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
	psqlconn := ""
	switch cfg.Driver {
	case "postgres":
		psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	case "mysql":
		psqlconn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, psqlconn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}
