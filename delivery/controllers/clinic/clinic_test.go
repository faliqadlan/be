package clinic

import (
	"be/entities"
	"errors"
	"testing"
)

type mockSuccess struct{}

func (m *mockSuccess) Create(clinicReq entities.Clinic) (entities.Clinic, error) {
	return entities.Clinic{}, nil
}

func (m *mockSuccess) Update(clinic_uid string, up entities.Clinic) (entities.Clinic, error) {
	return entities.Clinic{}, nil
}

func (m *mockSuccess) Delete(clinic_uid string) (entities.Clinic, error) {
	return entities.Clinic{}, nil
}

type mockFail struct{}

func (m *mockFail) Create(clinicReq entities.Clinic) (entities.Clinic, error) {
	return entities.Clinic{}, errors.New("")
}

func (m *mockFail) Update(clinic_uid string, up entities.Clinic) (entities.Clinic, error) {
	return entities.Clinic{}, errors.New("")
}

func (m *mockFail) Delete(clinic_uid string) (entities.Clinic, error) {
	return entities.Clinic{}, errors.New("")
}

func TestCreate(t *testing.T) {
	t.Run("success Create", func(t *testing.T) {
		
	})
}
