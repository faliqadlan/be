package patient

import (
	"be/configs"
	"be/delivery/controllers/auth"
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

func (m *mockSuccess) Create(patientReq entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *mockSuccess) Update(patient_uid string, req entities.Patient) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *mockSuccess) Delete(patient_uid string) (entities.Patient, error) {
	return entities.Patient{}, nil
}

func (m *mockSuccess) GetProfile(patient_uid string) (patient.Profile, error) {
	return patient.Profile{Image: "https://karen-givi-bucket.s3.ap-southeast-1.amazonaws.com/testing"}, nil
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

func (m *mockFail) GetProfile(patient_uid string) (patient.Profile, error) {
	return patient.Profile{}, errors.New("")
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
			"nik":      "123",
			"dob":      "05-05-2000",
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
			"password": 123,
			"nik":      "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validator username ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "",
			"email":    "doctor@",
			"password": "password",
			"nik":      "nik",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator email ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "userName",
			"email":    "",
			"password": "password",
			"nik":      "nik",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator password ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "userName",
			"email":    "user",
			"password": "",
			"nik":      "nik",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator nik ", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "userName",
			"email":    "user",
			"password": "password",
			"nik":      "",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
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
		writer.WriteField("nik", "nik")

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

	t.Run("error upload file", func(t *testing.T) {

		var reqBody = new(bytes.Buffer)

		var writer = multipart.NewWriter(reqBody)
		writer.WriteField("userName", "doctor1")
		writer.WriteField("email", "doctor")
		writer.WriteField("password", "doctor")
		writer.WriteField("nik", "nik")

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
			"password": "doctor",
			"nik":      "123",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &mockTaskS3M{})
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





		var e = echo.New()

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

		var controller = New(&mockSuccess{}, &mockTaskS3M{})
		controller.GetProfile()(context)

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
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, &mockTaskS3M{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
