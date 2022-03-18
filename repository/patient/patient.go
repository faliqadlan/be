package patient

import (
	"be/entities"
	"be/utils"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(req entities.Patient) (entities.Patient, error) {

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Patient{}, errors.New("user name is already exist")
	}

	// check email

	var checkEmail = r.db.Model(&entities.Patient{}).Where("email = ?", req.Email).Select("user_name as UserName").Scan(&userNameCheck{})

	if checkEmail.RowsAffected != 0 {
		return entities.Patient{}, errors.New("email is already exist")
	}

	var uid string
	for {
		uid = strconv.Itoa(int(uuid.New().ID()))
		var find = entities.Patient{}
		var res = r.db.Model(&entities.Patient{}).Where("patient_uid = ?", uid).Find(&find)
		if res.RowsAffected == 0 {
			break
		}
	}
	var err error
	req.Password, err = utils.HashPassword(req.Password)
	if err != nil {
		return entities.Patient{}, errors.New("error in hash password")
	}
	req.Patient_uid = uid

	if res := r.db.Model(&entities.Patient{}).Create(&req); res.Error != nil {
		return entities.Patient{}, res.Error
	}

	return req, nil
}

func (r *Repo) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {

	var resInit entities.Patient

	// check username

	type userNameCheck struct {
		UserName string
	}

	var checkUserName = r.db.Raw("? union all ? ", r.db.Model(&entities.Patient{}).Select("user_name").Where("user_name = ?", req.UserName), r.db.Model(&entities.Doctor{}).Select("user_name").Where("user_name = ?", req.UserName)).Scan(&userNameCheck{})

	if checkUserName.RowsAffected != 0 {
		return entities.Patient{}, errors.New("user name is already exist")
	}

	// check email

	var checkEmail = r.db.Model(&entities.Patient{}).Where("email = ?", req.Email).Select("user_name as UserName").Scan(&userNameCheck{})

	if checkEmail.RowsAffected != 0 {
		return entities.Patient{}, errors.New("email is already exist")
	}

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Find(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Patient{}, gorm.ErrRecordNotFound
	}

	var timeInit time.Time
	if resInit.CreatedAt != resInit.UpdatedAt {
		switch {
		case req.Nik != "":
			return entities.Patient{}, errors.New("nik can't updated")
		case req.PlaceBirth != "":
			return entities.Patient{}, errors.New("place birth can't updated")
		case req.Dob != datatypes.Date(timeInit):
			return entities.Patient{}, errors.New("date of birth can't updated")

		}
	}

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Updates(entities.Patient{
		UserName:   req.UserName,
		Email:      req.Email,
		Password:   req.Password,
		Nik:        req.Nik,
		Name:       req.Name,
		Image:      req.Image,
		Gender:     req.Gender,
		Address:    req.Address,
		PlaceBirth: req.PlaceBirth,
		Dob:        req.Dob,
		Job:        req.Job,
		Status:     req.Status,
		Religion:   req.Religion}); res.Error != nil || res.RowsAffected == 0 {
		switch {
		case res.Error == nil:
			return entities.Patient{}, gorm.ErrRecordNotFound
		default:
			return entities.Patient{}, res.Error
		}

	}

	return resInit, nil
}

func (r *Repo) Delete(patient_uid string) (entities.Patient, error) {
	var resInit entities.Patient

	if res := r.db.Model(&entities.Patient{}).Where("patient_uid = ?", patient_uid).Delete(&resInit); res.Error != nil || res.RowsAffected == 0 {
		return entities.Patient{}, gorm.ErrRecordNotFound
	}

	return resInit, nil
}

func (r *Repo) GetProfile(patient_uid, userName, email string) (Profile, error) {

	switch {
	case patient_uid != "":
		patient_uid = "patient_uid = '" + patient_uid + "'"
	case userName != "":
		patient_uid = "user_name = '" + userName + "'"
	case email != "":
		patient_uid = "email = '" + email + "'"
	}

	var profileResp Profile

	if res := r.db.Model(&entities.Patient{}).Where(patient_uid).Select("patient_uid as Patient_uid, nik as Nik, name as Name, image as Image, gender as Gender, address as Address, place_birth as PlaceBirth, date_format(dob, '%d-%m-%Y') as Dob, religion as Religion, status as Status, job as Job").Find(&profileResp); res.Error != nil || res.RowsAffected == 0 {
		return Profile{}, gorm.ErrRecordNotFound
	}

	return profileResp, nil
}
