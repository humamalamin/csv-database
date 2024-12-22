package worker

import (
	"database/sql"
	"sync"
)

func DispatchWorkers(db *sql.DB, jobs <-chan []interface{}, totalWorker int, process func(workerIndex int, db *sql.DB, job []interface{})) {
	var wg sync.WaitGroup

	for workerIndex := 0; workerIndex < totalWorker; workerIndex++ {
		wg.Add(1)
		go func(workerIndex int) {
			defer wg.Done()

			for job := range jobs {
				process(workerIndex, db, job)
			}
		}(workerIndex)
	}

	wg.Wait()
}
