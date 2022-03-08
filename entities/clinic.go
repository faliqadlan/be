package entities

import (
	"time"

	"gorm.io/gorm"
)

type Clinic struct {
	ID         uint `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Clinic_uid string         `gorm:"index;type:varchar(22)"`
	UserName   string         `gorm:"index;not null;type:varchar(100)"`
	Email      string         `gorm:"index;not null;type:varchar(100)"`
	Password   string         `gorm:"not null;type:varchar(100)"`
	DocterName string
	ClinicName string
	Address    string
	OpenDay    string `gorm:"type:enum('senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu');default:'senin'"`
	CloseDay   string `gorm:"type:enum('senin', 'selasa', 'rabu', 'kamis', 'jumat', 'sabtu', 'minggu');default:'senin'"`
	Capacity   int
	Doctors    []Doctor `gorm:"foreignKey:Clinic_uid;references:Clinic_uid"`
	Visits     []Visit  `gorm:"foreignKey:Clinic_uid;references:Clinic_uid"`
}
