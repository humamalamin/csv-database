package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func InsertBatch(db *sql.DB, tableName string, headers []string, batch [][]interface{}, driver string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(headers, ","),
		strings.Join(generateQuestionsMark(len(headers), driver), ","),
	)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, values := range batch {
		if _, err := stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return tx.Commit()
}

func generateQuestionsMark(n int, driver string) []string {
	s := make([]string, n)
	for i := 1; i <= n; i++ {
		if driver == "postgres" {
			s[i-1] = fmt.Sprintf("$%d", i)
		} else if driver == "mysql" {
			s[i-1] = "?"
		}

	}
	return s
}
