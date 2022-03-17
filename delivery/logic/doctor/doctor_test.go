package doctor

import (
	"be/entities"
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestValidationRequest(t *testing.T) {
	t.Run("error user name", func(t *testing.T) {
		var req Req

		req.UserName = "hotaru123   "

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})
	t.Run("error name", func(t *testing.T) {
		var req Req

		req.Name = "hotaru123   "

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error address", func(t *testing.T) {
		var req Req

		req.Address = "hotaru123   *(&*&**&*&*&(*&(*&&^"

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error status", func(t *testing.T) {
		var req Req

		req.Status = "123456789123456a"

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error open day", func(t *testing.T) {
		var req Req

		req.OpenDay = "123456789123456a"

		var l = New()

		var res = l.ValidationRequest(req)

		assert.NotNil(t, res)
		log.Info(res)
	})

	t.Run("error close day", func(t *testing.T) {
		var req Req

		req.CloseDay = "123456789123456a"

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("succeess to entity", func(t *testing.T) {
		var req Req

		res := req.ToDoctor()

		assert.Equal(t,&entities.Doctor{} ,res)
		log.Info(res)
	})

	t.Run("success", func(t *testing.T) {
		var req Req

		var l = New()

		var err = l.ValidationRequest(req)

		assert.Nil(t, err)
		log.Info(err)
	})
}
