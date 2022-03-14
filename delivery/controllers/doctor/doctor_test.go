package doctor

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"bytes"
	"encoding/json"
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

func TestCreate(t *testing.T) {
	t.Run("success Create", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			var e = echo.New()

			var reqBody, _ = json.Marshal(map[string]interface{}{
				"userName": "doctor1",
				"email":    "doctor@",
				"password": "a",
			})

			var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			var res = httptest.NewRecorder()
			req.Header.Set("Content-Type", "application/json")

			context := e.NewContext(req, res)
			context.SetPath("/doctor")

			var controller = New(&mockSuccess{}, sess)
			controller.Create()(context)

			var response = ResponseFormat{}

			json.Unmarshal([]byte(res.Body.Bytes()), &response)

			assert.Equal(t, 201, response.Code)
		})
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": 123,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, sess)
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validator", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"userName": "doctor1",
			"email":    "doctor@",
			"password": "",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
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
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockFail{}, sess)
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
			"password": "doctor",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createCapacity{}, sess)
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
			"password": "doctor",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&createUserName{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Update())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("file", func(t *testing.T) {

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

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockFail{}, sess)
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

		var controller = New(&createUserName{}, sess)
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

		var controller = New(&createCapacity{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockFail{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockFail{}, sess)
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

		var controller = New(&mockSuccess{}, sess)
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

		var controller = New(&mockFail{}, sess)
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetAll())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
