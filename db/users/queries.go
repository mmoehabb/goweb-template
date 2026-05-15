package users

import (
	"errors"

	"goweb/db"
)

func Add(username, password string) error {
	var existing DataModel
	result := db.Connection.Where("username = ?", username).First(&existing)
	if result.RowsAffected > 0 {
		return errors.New("username already found")
	}

	newUser := DataModel{Username: username, Password: password}
	result = db.Connection.Create(&newUser)
	return result.Error
}

func Get(username string) (DataModel, error) {
	var user DataModel
	result := db.Connection.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return DataModel{}, errors.New("couldn't find username")
	}
	return user, nil
}