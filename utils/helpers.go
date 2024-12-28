package utils

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/humamalamin/csv-database/csv"
	dbRepo "github.com/humamalamin/csv-database/db"
	"github.com/humamalamin/csv-database/worker"
	"github.com/jackc/pgx/v5/pgxpool"
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

	switch dbConfig.Driver {
	case "postgres":
		conectDB, err := dbRepo.OpenConnectPostgres(dbConfig)
		if err != nil {
			return err
		}

		dbConn, err := conectDB.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Failed to acquire connection: %v", err)
		}
		defer dbConn.Release()

		go worker.DispatchWorkersPsql(dbConn, jobs, worker.WorkerConfig{
			WorkerCount: 10,
			BatchSize:   500,
		}, func(workerIndex int, db *pgxpool.Conn, batch [][]interface{}) {
			if err := dbRepo.BulkInsertWithTransactionPsql(dbConn, tableName, headers, batch, dbConfig.Driver, 10000); err != nil {
				log.Printf("Worker %d failed: %v", workerIndex, err)
			}
		})
	case "mysql":
		conectDB, err := dbRepo.OpenConnectMysql(dbConfig)
		if err != nil {
			return err
		}

		defer conectDB.Close()

		go worker.DispatchWorkersMysql(conectDB, jobs, worker.WorkerConfig{
			WorkerCount: 10,
			BatchSize:   500,
		}, func(workerIndex int, db *sql.DB, batch [][]interface{}) {
			if err := dbRepo.BulkInsertWithTransactionMysql(db, tableName, headers, batch, dbConfig.Driver, 10000); err != nil {
				log.Printf("Worker %d failed: %v", workerIndex, err)
			}
		})
	}

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
