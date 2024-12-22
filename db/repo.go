package db

import (
	"database/sql"
	"fmt"
	"strings"
)

func InsertBatch(db *sql.DB, tableName string, headers []string, batch [][]interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(headers, ","),
		strings.Join(generateQuestionsMark(len(headers)), ","),
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

func generateQuestionsMark(n int) []string {
	s := make([]string, n)
	for i := 1; i <= n; i++ {
		s[i-1] = fmt.Sprintf("$%d", i)
	}
	return s
}