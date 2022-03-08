package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}

type MockFailAuthLib struct{}

func (m *MockFailAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{}, errors.New("")
}

func TestLogin(t *testing.T) {
	t.Run("error in request for login user", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"email": "anonim@123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := New(&MockAuthLib{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 400, resp.Code)
	})

	t.Run("error internal server error for login user", func(t *testing.T) {
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

		authCont := New(&MockFailAuthLib{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 500, resp.Code)
	})

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

		authCont := New(&MockAuthLib{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 200, resp.Code)
		// log.Info(resp)
	})

}
