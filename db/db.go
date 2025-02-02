package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func connect() (*pgxpool.Conn, error) {
  var err error
  if pool == nil {
    pool, err = pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
  }
	if err != nil {
    return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	return pool.Acquire(context.Background())
}

// execute each query in queries slice without returning any results
// if it didn't return error, then all queries passed successfully.
func Queries(queries []string) error {
	conn, err := connect()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v", p)
		}
	}()
	defer conn.Release()
	for _, query := range queries {
		rows, err := conn.Query(context.Background(), query)
		if err != nil {
			return fmt.Errorf("%q\nexecuting query: %s", err, query)
		}
		for rows.Next() {
		}
		err = rows.Err()
		if err != nil {
			return err
		}
	}
	return nil
}

type Connection struct {
  conn *pgxpool.Conn
}

func GetConnection() (*Connection, error) {
  conn, err := connect();
  if err != nil {
    return nil, fmt.Errorf("Unable to get Connection: %v\n", err)
  }
  return &Connection{conn: conn}, err
}

// Release the connection; the user will not be able to query with this struct any more.
func (c *Connection) Close() {
	c.conn.Release()
}

// execute single query and return the result
func (c *Connection) Query(query string, args ...any) ([]any, error) {
	defer c.conn.Release()
	return c.SeqQuery(query, args...)
}

// just like Query but does not close the connection
// make sure to call Query in the last of the "sequence"
// or manually call Disconnect
func (c *Connection) SeqQuery(query string, args ...any) ([]any, error) {
  var conn = c.conn
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v", p)
		}
	}()
	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	var res = []any{}
	for rows.Next() {
		r, err := rows.Values()
		if err != nil {
			return res, err
		}
		res = append(res, r)
	}
	rows.Close()
	return res, err
}

// hardcoded sql queries to seed (initialize) the database placed here
func Seed() error {
	err := Queries([]string{
		"CREATE TABLE IF NOT EXISTS users (username VARCHAR(45) PRIMARY KEY, password VARCHAR(45) NOT NULL);",
	})
	return err
}
