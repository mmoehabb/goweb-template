package db

import (
  "fmt"
  "os"
  "context"
  "github.com/jackc/pgx/v5"
)

// this global var allows us to implement SeqQuery function
var conn *pgx.Conn

func connect() (*pgx.Conn, error) {
  // urlExample := "postgres://username:password@localhost:5432/database_name"
  var err error
  if (conn == nil) {
    conn, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
  }
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return conn, err
	}
  return conn, nil
}

func Disconnect() error {
  if conn == nil {
    return nil
  }
  err := conn.Close(context.Background())
  if err == nil {
    conn = nil
  }
  return err
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
  defer Disconnect()

  for _, query := range queries {
    _, err := conn.Query(context.Background(), query)
    if err != nil {
      return err
    }
  }
  return nil
}

// execute single query and return the result
func Query(query string, args ...any) ([]any, error) {
  conn, err := connect()
  if err != nil {
    return nil, err
  }
  defer func() {
    if p := recover(); p != nil {
      fmt.Fprintf(os.Stderr, "internal error: %v", p)
    }
  }()
  defer Disconnect()

  rows, err := conn.Query(context.Background(), query, args...)
  if err != nil {
    return nil, err
  }
  if rows.Next() {
    return rows.Values()
  } else {
    return []any{}, err
  }
}

// just like Query by does not close the connection
// make sure to call Query in the last of the "sequence"
// or manually call Disconnect
func SeqQuery(query string, args ...any) ([]any, error) {
  conn, err := connect()
  if err != nil {
    return nil, err
  }
  defer func() {
    if p := recover(); p != nil {
      fmt.Fprintf(os.Stderr, "internal error: %v", p)
    }
  }()

  rows, err := conn.Query(context.Background(), query, args...)
  if err != nil {
    return nil, err
  }
  if rows.Next() {
    return rows.Values()
  } else {
    return []any{}, err
  }
}

// hardcoded sql queries to seed (initialize) the database placed here
func Seed() error {
  err := Queries([]string{ 
    "CREATE TABLE IF NOT EXISTS users (username VARCHAR(45) PRIMARY KEY, password VARCHAR(45) NOT NULL);",
  })
  return err
}

