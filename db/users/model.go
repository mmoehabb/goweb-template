package users

import "gorm.io/gorm"

type DataModel struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Nickname string
}

// TableName overrides the default table name (data_model; which is generated
// from the type name) to "users"
func (DataModel) TableName() string {
	return "users"
}
