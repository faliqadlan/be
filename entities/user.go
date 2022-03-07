package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	User_uid string `gorm:"index;type:varchar(22)"`
	Name     string `gorm:"not null;type:varchar(100)"`
	Email    string `gorm:"index;not null;type:varchar(100)"`
	Password string `gorm:"not null;type:varchar(100)"`
}
