package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"goweb/constants"
)

var pool interface{} // kept for future use

func connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", constants.AppConfig.DatabaseUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	return db, nil
}

// execute each query in queries slice without returning any results
// if it didn't return error, then all queries passed successfully.
func Queries(queries []string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("%q\nexecuting query: %s", err, query)
		}
	}
	return nil
}

type Connection struct {
	db *sql.DB
}

func GetConnection() (*Connection, error) {
	db, err := connect()
	if err != nil {
		return nil, fmt.Errorf("Unable to get Connection: %v\n", err)
	}
	return &Connection{db: db}, err
}

// Release the connection; the user will not be able to query with this struct any more.
func (c *Connection) Close() {
	c.db.Close()
}

// execute single query and return the result
func (c *Connection) Query(query string, args ...any) ([]any, error) {
	defer c.db.Close()
	return c.SeqQuery(query, args...)
}

// just like Query but does not close the connection
// make sure to call Query in the last of the "sequence"
// or manually call Disconnect
func (c *Connection) SeqQuery(query string, args ...any) ([]any, error) {
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var res []any
	for rows.Next() {
		values := make([]any, len(cols))
		for i := range values {
			var v any
			values[i] = &v
		}
		if err := rows.Scan(values...); err != nil {
			return res, err
		}
		rowData := make([]any, len(values))
		for i, v := range values {
			if b, ok := v.(*any); ok && b != nil {
				rowData[i] = *b
			} else {
				rowData[i] = v
			}
		}
		res = append(res, rowData)
	}
	return res, rows.Err()
}

// hardcoded sql queries to seed (initialize) the database placed here
func Seed() error {
	return nil
}

func RunMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	db, err := sql.Open("pgx", constants.AppConfig.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("failed to open db connection: %w", err)
	}
	defer db.Close()

	migrationsDir := filepath.Join(wd, "db", "migrations")
	goose.SetTableName("goose_db_version")
	if err := goose.Run("up", db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}
