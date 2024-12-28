package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BulkInsertWithTransactionPsql(
	db *pgxpool.Conn,
	table string,
	headers []string,
	data [][]interface{},
	driver string,
	batchSize int,
) error {
	ctx := context.Background()

	// Mulai transaksi
	tx, err := db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	// Panggil fungsi untuk batch insert
	err = InsertBatchPsql(tx, table, headers, data, batchSize)
	if err != nil {
		_ = tx.Rollback(ctx)
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func BulkInsertWithTransactionMysql(db *sql.DB, table string, headers []string, data [][]interface{}, driver string, batchSize int) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	err = InsertBatchMysql(tx, table, headers, data, driver, batchSize)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert batch: %w", err)
	}

	return tx.Commit()
}

func InsertBatchPsql(tx pgx.Tx, tableName string, headers []string, batch [][]interface{}, batchSize int) error {
	rows := len(batch)
	if rows == 0 {
		return fmt.Errorf("no data to insert")
	}

	cols := len(headers)
	if cols == 0 {
		return fmt.Errorf("no columns specified")
	}

	for i := 0; i < rows; i += batchSize {
		end := i + batchSize
		if end > rows {
			end = rows
		}

		currentBatch := batch[i:end]
		placeholders, args := generatePgxPlaceholdersAndValues(len(currentBatch), cols, currentBatch)

		query := fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES %s",
			tableName,
			strings.Join(headers, ", "),
			placeholders,
		)

		// Execute query within the transaction
		_, err := tx.Exec(context.Background(), query, args...)
		if err != nil {
			return fmt.Errorf("failed batch insert: %w", err)
		}
	}

	return nil
}

func InsertBatchMysql(tx *sql.Tx, tableName string, headers []string, batch [][]interface{}, driver string, batchSize int) error {
	rows := len(batch)
	if rows == 0 {
		return fmt.Errorf("no data to insert")
	}

	cols := len(headers)
	if cols == 0 {
		return fmt.Errorf("no columns specified")
	}

	for i := 0; i < rows; i += batchSize {
		end := i + batchSize
		if end > rows {
			end = rows
		}

		batch := batch[i:end]
		placeholders, values := generateBulkPlaceholders(len(batch), cols, driver)
		for _, row := range batch {
			values = append(values, row...)
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Join(headers, ", "), placeholders)
		_, err := tx.Exec(query, values...)
		if err != nil {
			return fmt.Errorf("failed batch insert: %w", err)
		}
	}
	return nil
}

func generateBulkPlaceholders(rows, cols int, driver string) (string, []interface{}) {
	var placeholders []string
	var values []interface{}

	for i := 0; i < rows; i++ {
		rowPlaceholders := make([]string, cols)
		for j := 0; j < cols; j++ {
			if driver == "postgres" {
				rowPlaceholders[j] = fmt.Sprintf("$%d", i*cols+j+1)
			} else {
				rowPlaceholders[j] = "?"
			}
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	return strings.Join(placeholders, ", "), values
}

func generatePgxPlaceholdersAndValues(batchCount, colCount int, batch [][]interface{}) (string, []interface{}) {
	var placeholders []string
	var args []interface{}

	for _, row := range batch {
		var rowPlaceholders []string
		for j := range row {
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", len(args)+j+1))
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
		args = append(args, row...)
	}

	return strings.Join(placeholders, ", "), args
}
