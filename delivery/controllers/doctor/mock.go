package doctor

import (
	"be/entities"
	"be/repository/doctor"
	"errors"
	"mime/multipart"
)

type mockTaskS3M struct{}

func (m *mockTaskS3M) UploadFileToS3(fileHeader multipart.FileHeader) (string, error) {
	return "", nil
}

func (m *mockTaskS3M) UpdateFileS3(name string, fileHeader multipart.FileHeader) string {
	return "success"
}

func (m *mockTaskS3M) DeleteFileS3(name string) string {
	return "success"
}


type mockSuccess struct{}

func (m *mockSuccess) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *mockSuccess) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *mockSuccess) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *mockSuccess) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{}, nil
}

func (m *mockSuccess) GetAll() (doctor.All, error) {
	return doctor.All{}, nil
}

type mockFail struct{}

func (m *mockFail) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("")
}

func (m *mockFail) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("")
}

func (m *mockFail) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("")
}

func (m *mockFail) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{}, errors.New("")
}

func (m *mockFail) GetAll() (doctor.All, error) {
	return doctor.All{}, errors.New("")
}

type createCapacity struct{}

func (m *createCapacity) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("can't assign capacity below zero")
}

func (m *createCapacity) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("can't update capacity below total pending patients")
}

func (m *createCapacity) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("")
}

func (m *createCapacity) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{}, errors.New("")
}

func (m *createCapacity) GetAll() (doctor.All, error) {
	return doctor.All{}, errors.New("")
}

type createUserName struct{}

func (m *createUserName) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("user name is already exist")
}

func (m *createUserName) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("user name is already exist")
}

func (m *createUserName) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, errors.New("")
}

func (m *createUserName) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{}, errors.New("")
}

func (m *createUserName) GetAll() (doctor.All, error) {
	return doctor.All{}, errors.New("")
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}
