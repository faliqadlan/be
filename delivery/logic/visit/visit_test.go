package visit

import (
	"testing"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestToVIsit(t *testing.T) {

	t.Run("succeess", func(t *testing.T) {

		var req Req

		req.Status = "pending"

		var l = New()

		err := l.ValidationRequest(req)

		assert.Nil(t, err)
		log.Info(err)
	})

	t.Run("error empyty", func(t *testing.T) {

		var req Req

		var l = New()

		err := l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error status", func(t *testing.T) {

		var req Req

		req.Status = "status"

		var l = New()

		err := l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error parsing date", func(t *testing.T) {

		var req Req

		req.Date = "05-05-00"

		_, err := req.ToVisit()

		assert.NotNil(t, err)
		log.Info(err)
	})

	// t.Run("error past date", func(t *testing.T) {

	// 	var req Req

	// 	var layout = "02-01-2006"

	// 	req.Date = time.Now().AddDate(0,0,-1).Format(layout)

	// 	_, err := req.ToVisit()

	// 	assert.NotNil(t, err)
	// 	log.Info(err)
	// })

	t.Run("succeess", func(t *testing.T) {

		var req Req

		var layout = "02-01-2006"

		req.Date = time.Now().Format(layout)

		_, err := req.ToVisit()

		assert.Nil(t, err)
		log.Info(err)
	})
}
