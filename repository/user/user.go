package user

import (
	"be/entities"
	"be/utils"
	"errors"

	"github.com/lithammer/shortuuid"
	"gorm.io/gorm"
)

type UserDb struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserDb {
	return &UserDb{
		db: db,
	}
}

func (repo *UserDb) Create(user entities.User) (entities.User, error) {

	userInit := entities.User{}

	checkEmail := repo.db.Where("email = ?", user.Email).Find(&userInit)

	if checkEmail.RowsAffected != 0 {
		return entities.User{}, errors.New("email already exist")
	}

	var uid string

	for {
		uid = shortuuid.New()
		userFind := entities.User{}
		res := repo.db.Model(&entities.User{}).Where("user_uid = ?", uid).First(&userFind)
		if res.RowsAffected == 0 {
			break
		}
	}
	user.Password, _ = utils.HashPassword(user.Password)

	user.User_uid = uid

	if err := repo.db.Create(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user, nil
}
