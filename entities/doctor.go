package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Doctor struct {
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	Doctor_uid          string         `gorm:"index;type:varchar(22);primaryKey"`
	Clinic_uid          string         `gorm:"index;type:varchar(22)"`
	UserName            string         `gorm:"index;not null;type:varchar(100)"`
	Email               string         `gorm:"index;not null;type:varchar(100)"`
	Password            string         `gorm:"not null;type:varchar(100)"`
	Nik                 string         `gorm:"type:varchar(16)"`
	FullName            string
	Address             string `gorm:"not null"`
	PlaceBirth          string `gorm:"type:varchar(100)"`
	Dob                 datatypes.Date
	NoDegreeCertificate string
	Status              string   `gorm:"type:enum('belumKawin', 'kawin', 'ceraiHidup', 'ceraiMati', 'lainnya');default:'lainnya'"`
	Religion            string   `gorm:"type:enum('islam', 'krister', 'katolik', 'budha', 'hindu', 'konghuchu', 'lainnya');default:'lainnya'"`
	Visits             []Visit `gorm:"foreignKey:Doctor_uid;references:Doctor_uid"`
}
