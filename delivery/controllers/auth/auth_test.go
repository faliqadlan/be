package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

type MockAuthLibFailToken struct{}

func (m *MockAuthLibFailToken) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "",
		"type": "clinic",
	}, nil
}

type MockIncorrectPassword struct{}

func (m *MockIncorrectPassword) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{}, errors.New("incorrect password")
}

type DeletedAccount struct{}

func (m *DeletedAccount) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{}, errors.New("account is deleted")
}

type AccountNotFound struct{}

func (m *AccountNotFound) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{}, gorm.ErrRecordNotFound
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

	t.Run("there's some problem is server", func(t *testing.T) {
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

	t.Run("incorrect password", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "incorrect password",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := New(&MockIncorrectPassword{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 500, resp.Code)
	})

	t.Run("account is deleted", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "account is deleted",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := New(&DeletedAccount{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 500, resp.Code)
	})

	t.Run(gorm.ErrRecordNotFound.Error(), func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": gorm.ErrRecordNotFound.Error(),
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := New(&AccountNotFound{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 500, resp.Code)
	})

	t.Run("error in process token", func(t *testing.T) {
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

		authCont := New(&MockAuthLibFailToken{})
		authCont.Login()(context)

		resp := LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)
		assert.Equal(t, 406, resp.Code)
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
		log.Info(resp)
	})

}
