package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Patient struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Patient_uid string         `gorm:"index;type:varchar(22)"`
	Name        string         `gorm:"not null;type:varchar(100)"`
	Email       string         `gorm:"index;not null;type:varchar(100)"`
	Password    string         `gorm:"not null;type:varchar(100)"`
	Nik         int            `gorm:"type:DECIMAL(16)"`
	Address     string         `gorm:"not null"`
	PlaceBirth  string         `gorm:"type:varchar(100)"`
	Dob         datatypes.Date
	Job         string
	Status      string
	Religion    string
}
