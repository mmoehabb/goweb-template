package db

import (
  "fmt"
  "os"
  "context"
  "github.com/jackc/pgx/v5"
)

// execute each query in queries slice without returning any results
// if it didn't return error, then all queries passed successfully.
func Query(queries []string) error {
  conn, err := connect()
  if err != nil {
    return err
  }
  defer func() {
    if p := recover(); p != nil {
      fmt.Fprintf(os.Stderr, "internal error: %v", p)
    }
  }()
  defer disconnect(conn)

  for _, query := range queries {
    _, err := conn.Query(context.Background(), query)
    if err != nil {
      return err
    }
  }
  return nil
}

// execute single query and return the result
func QueryResult(query string, args ...any) ([]any, error) {
  conn, err := connect()
  if err != nil {
    return nil, err
  }
  defer func() {
    if p := recover(); p != nil {
      fmt.Fprintf(os.Stderr, "internal error: %v", p)
    }
  }()
  defer disconnect(conn)

  rows, err := conn.Query(context.Background(), query, args...)
  if err != nil {
    return nil, err
  }
  if rows.Next() {
    return rows.Values()
  } else {
    return []any{}, nil
  }
}

func connect() (*pgx.Conn, error) {
  // urlExample := "postgres://username:password@localhost:5432/database_name"
  conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return conn, err
	}
  return conn, nil
}

func disconnect(conn *pgx.Conn) error {
  return conn.Close(context.Background())
}
