package db

import (
	"gorm.io/gorm"
)

var originalConnection *gorm.DB

func SetConnectionForTest(conn *gorm.DB) {
	originalConnection = Connection
	Connection = conn
}

func ResetConnectionForTest() {
	if originalConnection != nil {
		Connection = originalConnection
	}
}
