package user

import (
	"errors"
	"goweb/db"
)

type User struct{
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
    return err
  }
  return nil
}

func Get(username string) (User, error) {
  res, err := db.Query("SELECT * FROM users WHERE username=$1", username)
  if err != nil {
    return User{}, err
  }
  if len(res) == 0 {
    return User{}, errors.New("couldn't find username.")
  }
  user := User{ 
    Username: res[0].(string), 
    Password: res[1].(string),
  }
  return user, nil
}
