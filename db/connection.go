package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
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

func OpenDbConnection(cfg ConfigDB) (interface{}, error) {
	// log.Println("=> open db connection")
	// var connStr string

	// switch cfg.Driver {
	// case "postgres":
	// 	connStr = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
	// 		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	// 	return connectPostgres(connStr, cfg), nil
	// case "mysql":
	// 	connStr = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
	// 		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	// 	return connStr, nil
	// default:
	// 	return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	// }

	// db, err := sql.Open("pgx", connStr)
	// if err != nil {
	// 	return nil, err
	// }

	// db.SetMaxOpenConns(cfg.MaxConns)
	// db.SetMaxIdleConns(cfg.MaxIdleConns)

	// return db, nil
	panic("STOP")
}

func OpenConnectPostgres(cfg ConfigDB) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?connect_timeout=10&statement_timeout=5000",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Ensure the connection pool is healthy
	var countConn int32 = 0
	if conn.Stat().TotalConns() == int32(countConn) {
		log.Fatalf("no connections available in pool")
	}
	return conn, nil
}

func OpenConnectMysql(cfg ConfigDB) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return db, nil
}
