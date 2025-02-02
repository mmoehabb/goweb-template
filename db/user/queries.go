package user

import (
	"errors"
	"fmt"
	"goweb/db"
	"os"
)

type User struct {
	Username string
	Password string
}

func Add(username, password string) error {
  conn, err := db.GetConnection();
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Internal Server Error.")
	}

	res, err := conn.SeqQuery("SELECT * FROM users WHERE username=$1", username)
	if len(res) != 0 {
		conn.Close()
		return errors.New("username already found.")
	}
	_, err = conn.Query("INSERT INTO users VALUES ($1, $2)", username, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Internal Server Error.")
	}
	return nil
}

func Get(username string) (User, error) {
  conn, err := db.GetConnection();
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return User{}, errors.New("Internal Server Error.")
	}

	res, err := conn.Query("SELECT * FROM users WHERE username=$1", username)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return User{}, errors.New("Internal Server Error.")
	}
	if len(res) == 0 {
		return User{}, errors.New("couldn't find username.")
	}
	row := res[0].([]any)
	user := User{
		Username: row[0].(string),
		Password: row[1].(string),
	}
	return user, nil
}
