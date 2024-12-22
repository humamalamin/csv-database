package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/humamalamin/csv-database/csv"
	dbRepo "github.com/humamalamin/csv-database/db"
	"github.com/humamalamin/csv-database/worker"
	_ "github.com/lib/pq"
)

func main() {
	start := time.Now()

	// Konfigurasi database
	dbConfig := dbRepo.ConfigDB{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "lokal",
		DBName:       "csv_database",
		MaxConns:     20,
		MaxIdleConns: 10,
	}

	database, err := dbRepo.OpenDbConnection(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Buka file CSV
	csvReader, file, err := csv.OpenCsvFile("annual-enterprise-survey-2023-financial-year-provisional-size-bands.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Baca header
	headers, err := csv.ReadHeaders(csvReader)
	if err != nil {
		log.Fatal(err)
	}

	// Worker processing
	jobs := make(chan []interface{})
	go worker.DispatchWorkers(database, jobs, 10, func(workerIndex int, db *sql.DB, job []interface{}) {
		if err := dbRepo.InsertBatch(db, "annual_enterprise", headers, [][]interface{}{job}); err != nil {
			log.Printf("Worker %d failed: %v", workerIndex, err)
		}
	})

	// Kirim data ke worker
	for {
		row, err := csvReader.Read()
		if err != nil {
			break
		}
		jobs <- toInterfaceSlice(row)
	}
	close(jobs)

	duration := time.Since(start)
	log.Printf("Done in %s\n", duration)
}

func toInterfaceSlice(strSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
