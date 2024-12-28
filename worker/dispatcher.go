package worker

import (
	"database/sql"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerConfig struct {
	WorkerCount int
	BatchSize   int
}

func DispatchWorkersMysql(db *sql.DB, jobs <-chan []interface{}, cfg WorkerConfig, insertFunc func(workerIndex int, db *sql.DB, batch [][]interface{})) {
	var wg sync.WaitGroup

	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			var batch [][]interface{}
			for job := range jobs {
				batch = append(batch, job)
				if len(batch) >= cfg.BatchSize {
					insertFunc(workerIndex, db, batch)
					batch = nil // Reset batch
				}
			}
			// Process remaining batch
			if len(batch) > 0 {
				insertFunc(workerIndex, db, batch)
			}
		}(i)
	}
	wg.Wait()
	log.Println("All workers finished")
}

func DispatchWorkersPsql(db *pgxpool.Conn, jobs <-chan []interface{}, cfg WorkerConfig, insertFunc func(workerIndex int, db *pgxpool.Conn, batch [][]interface{})) {
	var wg sync.WaitGroup

	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()
			var batch [][]interface{}
			for job := range jobs {
				batch = append(batch, job)
				if len(batch) >= cfg.BatchSize {
					insertFunc(workerIndex, db, batch)
					batch = nil // Reset batch
				}
			}
			// Process remaining batch
			if len(batch) > 0 {
				insertFunc(workerIndex, db, batch)
			}
		}(i)
	}
	wg.Wait()
	log.Println("All workers finished")
}
