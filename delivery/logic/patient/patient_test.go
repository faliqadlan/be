package patient

import (
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestValidationStruct(t *testing.T) {
	t.Run("validator nik", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator name", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator gender", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator address", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator placeBirth", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "",
			Dob:        "dob",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator dob", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator job", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Job:        "",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator status", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Job:        "job",
			Status:     "",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator religion", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Job:        "job",
			Status:     "status",
			Religion:   "",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("validator succeess", func(t *testing.T) {
		var req = Req{
			UserName:   "userName",
			Email:      "email",
			Password:   "password",
			Nik:        "nik",
			Name:       "name",
			Gender:     "gender",
			Address:    "address",
			PlaceBirth: "placeBirth",
			Dob:        "dob",
			Job:        "job",
			Status:     "status",
			Religion:   "religion",
		}
		var l = New()
		var err = l.ValidationStruct(req)
		assert.Nil(t, err)
		log.Info(err)
	})
}

func TestValidationRegexPatient(t *testing.T) {
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

	t.Run("error nik", func(t *testing.T) {
		var req Req

		req.Nik = "123456789123456a"

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("error date", func(t *testing.T) {
		var req Req

		req.Dob = "123456789123456a"

		_, err := req.ToPatient()

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("succeess date", func(t *testing.T) {
		var req Req

		req.Dob = "05-05-2002"

		_, err := req.ToPatient()

		assert.Nil(t, err)
		log.Info(err)
	})

	t.Run("succeess no date", func(t *testing.T) {
		var req Req

		_, err := req.ToPatient()

		assert.Nil(t, err)
		log.Info(err)
	})

	t.Run("error gender", func(t *testing.T) {
		var req Req

		req.Gender = "123456789123456a"

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

	t.Run("error religion", func(t *testing.T) {
		var req Req

		req.Religion = "123456789123456a"

		var l = New()

		var err = l.ValidationRequest(req)

		assert.NotNil(t, err)
		log.Info(err)
	})

	t.Run("success", func(t *testing.T) {
		var req Req

		var l = New()

		var err = l.ValidationRequest(req)

		assert.Nil(t, err)
		log.Info(err)
	})
}
