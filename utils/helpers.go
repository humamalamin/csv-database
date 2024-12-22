package utils

import "log"

func LogProgress(workerIndex, counter, step int) {
	if counter%step == 0 {
		log.Printf("=> worker %d processed %d rows\n", workerIndex, counter)
	}
}
