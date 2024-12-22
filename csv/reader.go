package csv

import (
	"encoding/csv"
	"os"
)

func OpenCsvFile(filePath string) (*csv.Reader, *os.File, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(f)
	return reader, f, nil
}

func ReadHeaders(csvReader *csv.Reader) ([]string, error) {
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	return headers, nil
}
