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
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

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

func (m *mockSuccess) GetPatients(doctor_uid string) (doctor.PatientsResp, error) {
	return doctor.PatientsResp{}, nil
}

func (m *mockSuccess) GetDashboard(doctor_uid string) (doctor.Dashboard, error) {
	return doctor.Dashboard{}, nil
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

func (m *mockFail) GetPatients(doctor_uid string) (doctor.PatientsResp, error) {
	return doctor.PatientsResp{}, errors.New("")
}

func (m *mockFail) GetDashboard(doctor_uid string) (doctor.Dashboard, error) {
	return doctor.Dashboard{}, errors.New("")
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}

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

			var controller = New(&mockSuccess{})
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

		var controller = New(&mockSuccess{})
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

		var controller = New(&mockSuccess{})
		controller.Create()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
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

		var controller = New(&mockFail{})
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

		var controller = New(&mockSuccess{})
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

		var controller = New(&mockSuccess{})
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

		var controller = New(&mockFail{})
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

		var controller = New(&mockSuccess{})
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

		var controller = New(&mockFail{})
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

		var controller = New(&mockSuccess{})
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

		var controller = New(&mockFail{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetProfile())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}

func TestGetPatients(t *testing.T) {
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

		var controller = New(&mockSuccess{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetPatients())(context); err != nil {
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

		var controller = New(&mockFail{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetPatients())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}

func TestGetDashboard(t *testing.T) {
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

		var controller = New(&mockSuccess{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetDashboard())(context); err != nil {
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

		var controller = New(&mockFail{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetDashboard())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 500, response.Code)
	})

}
