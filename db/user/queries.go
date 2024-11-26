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
	res, err := db.SeqQuery("SELECT * FROM users WHERE username=$1", username)
	if len(res) != 0 {
		db.Disconnect()
		return errors.New("username already found.")
	}
	_, err = db.Query("INSERT INTO users VALUES ($1, $2)", username, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Internal Server Error.")
	}
	return nil
}

func Get(username string) (User, error) {
	res, err := db.Query("SELECT * FROM users WHERE username=$1", username)
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
