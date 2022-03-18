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
	Patient_uid string         `gorm:"index;type:varchar(22);primaryKey"`
	UserName    string         `gorm:"index;not null;type:varchar(100)"`
	Email       string         `gorm:"index;not null;type:varchar(100)"`
	Password    string         `gorm:"not null;type:varchar(100)"`
	Nik         string         `gorm:"type:varchar(16)"`
	Name        string
	Image       string `gorm:"default:'https://www.teralogistics.com/wp-content/uploads/2020/12/default.png'"`
	Gender      string `gorm:"type:enum('pria', 'wanita', 'lainnya');default:'lainnya'"`
	Address     string `gorm:"not null"`
	PlaceBirth  string `gorm:"type:varchar(100)"`
	Dob         datatypes.Date
	Job         string
	Status      string  `gorm:"type:enum('belumKawin', 'kawin', 'ceraiHidup', 'ceraiMati', 'lainnya');default:'lainnya'"`
	Religion    string  `gorm:"type:enum('islam', 'kristen', 'katolik', 'protestan', 'budha', 'hindu', 'konghuchu', 'lainnya');default:'lainnya'"`
	Visits      []Visit `gorm:"foreignKey:Patient_uid;references:Patient_uid"`
}
