package main

import (
	"log"

	dbRepo "github.com/humamalamin/csv-database/db"
	"github.com/humamalamin/csv-database/utils"
)

func main() {
	// Konfigurasi database
	dbConfig := dbRepo.ConfigDB{
		Driver:       "postgres",
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "lokal",
		DBName:       "csv_database",
		MaxConns:     20,
		MaxIdleConns: 10,
	}

	err := utils.ProcessCSVToDatabase(dbConfig, "annual-enterprise-survey-2023-financial-year-provisional-size-bands.csv", "annual_enterprise", 10)
	if err != nil {
		log.Fatal(err)
	}

}
