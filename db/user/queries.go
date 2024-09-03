package user

import "goweb/db"

type User struct{
  Username string
  Password string
}

func Add(username, password string) error {
  _, err := db.QueryResult("INSERT INTO users VALUES ($1, $2)", username, password)
  if err != nil {
    return err
  }
  return nil
}

func Get(username string) (User, error) {
  res, err := db.QueryResult("SELECT * FROM users WHERE username=$1", username)
  if err != nil {
    return User{}, err
  }
  user := User{ 
    Username: res[0].(string), 
    Password: res[1].(string),
  }
  return user, nil
}
