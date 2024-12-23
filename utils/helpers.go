package utils

import (
	"database/sql"
	"log"
	"time"

	"github.com/humamalamin/csv-database/csv"
	dbRepo "github.com/humamalamin/csv-database/db"
	"github.com/humamalamin/csv-database/worker"
)

func LogProgress(workerIndex, counter, step int) {
	if counter%step == 0 {
		log.Printf("=> worker %d processed %d rows\n", workerIndex, counter)
	}
}

func ProcessCSVToDatabase(
	dbConfig dbRepo.ConfigDB,
	csvFilePath string,
	tableName string,
	workerCount int,
) error {
	start := time.Now()

	// Buka koneksi database
	database, err := dbRepo.OpenDbConnection(dbConfig)
	if err != nil {
		return err
	}
	defer database.Close()

	// Buka file CSV
	csvReader, file, err := csv.OpenCsvFile(csvFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Baca header
	headers, err := csv.ReadHeaders(csvReader)
	if err != nil {
		return err
	}

	// Worker processing
	jobs := make(chan []interface{})
	go worker.DispatchWorkers(database, jobs, workerCount, func(workerIndex int, db *sql.DB, job []interface{}) {
		if err := dbRepo.InsertBatch(db, tableName, headers, [][]interface{}{job}, dbConfig.Driver); err != nil {
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
	log.Printf("Processed CSV to database in %s\n", duration)
	return nil
}

func toInterfaceSlice(strSlice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(strSlice))
	for i, v := range strSlice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
