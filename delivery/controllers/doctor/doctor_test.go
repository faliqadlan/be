package doctor

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/entities"
	"be/repository/doctor"
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
	return doctor.ProfileResp{Image: "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/testing"}, nil
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
	return doctor.ProfileResp{Image: "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png"}, nil
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

type recordNotFound struct{}

func (m *recordNotFound) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{}, gorm.ErrRecordNotFound
}

func (m *recordNotFound) GetAll() (doctor.All, error) {
	return doctor.All{}, gorm.ErrRecordNotFound
}

type updateFile struct{}

func (m *updateFile) Create(DoctorReq entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *updateFile) Update(Doctor_uid string, up entities.Doctor) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *updateFile) Delete(Doctor_uid string) (entities.Doctor, error) {
	return entities.Doctor{}, nil
}

func (m *updateFile) GetProfile(doctor_uid string) (doctor.ProfileResp, error) {
	return doctor.ProfileResp{Image: "https://www.teralogistics.com/wp-content/uploads/2020/12/default.png"}, nil
}

func (m *updateFile) GetAll() (doctor.All, error) {
	return doctor.All{}, nil
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}

func TestCreate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "a",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "a",
			"name":     50,
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validator userName", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "",
			"email":    "doctor@",
			"password": "a",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})
		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator email", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "",
			"password": "a",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator password", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")
		writer.WriteField("name", "name")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("status", "status")
		writer.WriteField("openDay", "senin")
		writer.WriteField("closeDay", "senin")
		writer.WriteField("capacity", "50")

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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
		// log.Info(response.Message)
	})

	t.Run("fail upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")
		writer.WriteField("name", "name")
		writer.WriteField("address", "123456789123456")
		writer.WriteField("status", "status")
		writer.WriteField("openDay", "senin")
		writer.WriteField("closeDay", "senin")
		writer.WriteField("capacity", "50")

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

		var controller = New(&mockSuccess{}, &failTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "abc",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("can't assign capacity below zero", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "abc",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createCapacity{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("user name is already exist", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "abc",
			"name":     "name",
			"address":  "123456789123456",
			"status":   "status",
			"openDay":  "openDay",
			"closeDay": "closeDay",
			"capacity": 50,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createUserName{}, &mockTaskS3M{})
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
			"name": "doctor name",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
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
			"name": 123,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("succeess upload", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("error get image link", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("error upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&mockSuccess{}, &failTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("update file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&updateFile{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("error update file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&createCapacity{}, &failTaskS3M{})
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
			"name": "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("user name is already exist", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"name": "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createUserName{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("can't update capacity below total pending patients", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"name": "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createCapacity{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("record not found", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"name": "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{})
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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
	})

	t.Run("fail delete file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")

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

		var controller = New(&mockSuccess{}, &failTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 202, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Delete())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("record not found", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&recordNotFound{}, &mockTaskS3M{})
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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor/profile")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor/profile")

		var controller = New(&recordNotFound{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
func TestGetAll(t *testing.T) {
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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetAll())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 200, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor/profile")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetAll())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
