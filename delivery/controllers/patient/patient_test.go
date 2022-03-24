package patient

import (
	"be/configs"
	"be/delivery/controllers/auth"
	logic "be/delivery/logic/patient"
	"be/entities"
	"be/repository/patient"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type successLogic struct{}

func (m *successLogic) ValidationRequest(req logic.Req) error {
	return nil
}

func (m *successLogic) ValidationStruct(req logic.Req) error {
	return nil
}

type errorLogic struct{}

func (m *errorLogic) ValidationRequest(req logic.Req) error {
	return errors.New("")
}

func (m *errorLogic) ValidationStruct(req logic.Req) error {
	return errors.New("")
}

type errorLogicStruct struct{}

func (m *errorLogicStruct) ValidationRequest(req logic.Req) error {
	return errors.New("")
}

func (m *errorLogicStruct) ValidationStruct(req logic.Req) error {
	return nil
}

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

type failTaskS3M struct{}

func (m *failTaskS3M) UploadFileToS3(fileHeader multipart.FileHeader) (string, error) {
	return "", errors.New("")
}

func (m *failTaskS3M) UpdateFileS3(name string, fileHeader multipart.FileHeader) string {
	return "error"
}

func (m *failTaskS3M) DeleteFileS3(name string) string {
	return "error"
}

type mockSuccess struct{}

func (m *mockSuccess) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{Patient_uid: "abcde"}, nil
}

func (m *mockSuccess) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *mockSuccess) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *mockSuccess) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{Image: "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/testing"}, nil
}

func (m *mockSuccess) GetAll() (patient.All, error) {
	return patient.All{}, nil
}

type defaultImage struct{}

func (m *defaultImage) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *defaultImage) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *defaultImage) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *defaultImage) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{Image: "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png"}, nil
}

func (m *defaultImage) GetAll() (patient.All, error) {
	return patient.All{}, nil
}

type mockFail struct{}

func (m *mockFail) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("")
}

func (m *mockFail) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("")
}

func (m *mockFail) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, errors.New("")
}

func (m *mockFail) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{}, errors.New("")
}

func (m *mockFail) GetAll() (patient.All, error) {
	return patient.All{}, errors.New("")
}

type recordNotFound struct{}

func (m *recordNotFound) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) GetAll() (patient.All, error) {
	return patient.All{}, nil
}

type userNameCheck struct{}

func (m *userNameCheck) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("user name is already exist")
}

func (m *userNameCheck) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("user name is already exist")
}

func (m *userNameCheck) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, errors.New("user name is already exist")
}

func (m *userNameCheck) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{}, errors.New("user name is already exist")
}

func (m *userNameCheck) GetAll() (patient.All, error) {
	return patient.All{}, nil
}

type emailCheck struct{}

func (m *emailCheck) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("email is already exist")
}

func (m *emailCheck) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, errors.New("email is already exist")
}

func (m *emailCheck) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, errors.New("user name is already exist")
}

func (m *emailCheck) GetProfile(patient_uid, userName, email string) (patient.Profile, error) {
	return patient.Profile{}, errors.New("user name is already exist")
}

func (m *emailCheck) GetAll() (patient.All, error) {
	return patient.All{}, nil
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"doctor_uid":"abcde",
		"type": "clinic",
	}, nil
}

func TestCreate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "a",
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   123,
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("binding username ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   50,
			"email":      "doctor@",
			"password":   "a",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding email ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      50,
			"password":   "a",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding password ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   50,
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding nik ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        50,
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding name ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       50,
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding gender ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     50,
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding address ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    50,
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding placeBirth ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": 50,
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding dob ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        50,
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding job ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        50,
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding status ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     50,
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("binding religion ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "status",
			"religion":   50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator struct ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "status",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &errorLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator request ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "status",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &errorLogicStruct{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("parsing dob", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "password",
			"nik":        "123",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "address",
			"placeBirth": "malang",
			"dob":        "05-05-00",
			"job":        "lainnya",
			"status":     "status",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("succeess upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &failTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("userNameCheck", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "a",
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&userNameCheck{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("email", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "a",
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&emailCheck{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName":   "doctor1",
			"email":      "doctor@",
			"password":   "a",
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}

func TestUpdate(t *testing.T) {
	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, 202, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       123,
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validation request", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "name",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &errorLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("parsing dob", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "name",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-00",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("succeess update file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("error get link file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &failTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("error update file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &failTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("succeess upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&defaultImage{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("fail upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("nik", "1234567891234567")
		writer.WriteField("name", "name")
		writer.WriteField("gender", "pria")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("placeBirth", "placeBirth")
		writer.WriteField("dob", "05-05-2002")
		writer.WriteField("job", "lainnya")
		writer.WriteField("status", "lainnya")
		writer.WriteField("religion", "religion")

		part, err := writer.CreateFormFile("file", "photo.jpg")
		if err != nil {
			log.Warn(err)
		}
		part.Write([]byte(`sample`))
		writer.Close()

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", reqBody)
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&defaultImage{}, &failTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("userNameCheck", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&userNameCheck{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("emailCheck", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&emailCheck{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"nik":        "1234567891234567",
			"name":       "subejo",
			"gender":     "pria",
			"address":    "123456789123456",
			"placeBirth": "malang",
			"dob":        "05-05-2000",
			"job":        "lainnya",
			"status":     "lainnya",
			"religion":   "lainnya",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}

func TestDelete(t *testing.T) {
	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("error delete file", func(t *testing.T) {

		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &failTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}

func TestGetProfile(t *testing.T) {
	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("success query param all", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/patient/profile")
		context.QueryParams().Add("all", "all")
		// log.Info(context.QueryString())

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.GetProfile()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("error query param all", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/patient/profile")
		context.QueryParams().Add("all", "all")
		// log.Info(context.QueryString())

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		controller.GetProfile()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("success query param", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/patient/profile")
		context.QueryParams().Add("patient_uid", "ini_query")
		// log.Info(context.QueryString())

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		controller.GetProfile()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}

func TestGetCheck(t *testing.T) {
	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetCheck())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetCheck())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{}, &successLogic{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetCheck())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
