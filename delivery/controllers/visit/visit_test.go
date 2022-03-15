package visit

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/entities"
	"be/repository/visit"
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
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm"
)

type mockSuccess struct{}

func (m *mockSuccess) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, nil
}

func (m *mockSuccess) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, nil
}

type errorVisitList struct{}

func (m *errorVisitList) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *errorVisitList) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *errorVisitList) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *errorVisitList) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, nil
}

func (m *errorVisitList) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, errors.New("")
}

type errorUpdateEventId struct{}

func (m *errorUpdateEventId) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *errorUpdateEventId) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *errorUpdateEventId) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *errorUpdateEventId) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, nil
}

func (m *errorUpdateEventId) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, nil
}

type mockFail struct{}

func (m *mockFail) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, errors.New("")
}

func (m *mockFail) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, errors.New("")
}

type spesificError struct{}

func (m *spesificError) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("there's another appoinment in pending")
}

func (m *spesificError) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *spesificError) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *spesificError) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, gorm.ErrRecordNotFound
}

func (m *spesificError) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, gorm.ErrRecordNotFound
}

type leftCapacity struct{}

func (m *leftCapacity) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("left capacity can't below zero")
}

func (m *leftCapacity) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *leftCapacity) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *leftCapacity) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, gorm.ErrRecordNotFound
}

func (m *leftCapacity) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, gorm.ErrRecordNotFound
}

type invalidDoctorUid struct{}

func (m *invalidDoctorUid) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("Cannot add or update a child row: a foreign key constraint fails (`crud_api_test`.`visits`, CONSTRAINT `fk_doctors_visits` FOREIGN KEY (`doctor_uid`) REFERENCES `doctors` (`doctor_uid`))")
}

func (m *invalidDoctorUid) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *invalidDoctorUid) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *invalidDoctorUid) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, gorm.ErrRecordNotFound
}

func (m *invalidDoctorUid) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, gorm.ErrRecordNotFound
}

type invalidPatientUid struct{}

func (m *invalidPatientUid) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("Cannot add or update a child row: a foreign key constraint fails (`crud_api_test`.`visits`, CONSTRAINT `fk_patients_visits` FOREIGN KEY (`patient_uid`) REFERENCES `patients` (`patient_uid`))")
}

func (m *invalidPatientUid) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *invalidPatientUid) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, gorm.ErrRecordNotFound
}

func (m *invalidPatientUid) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, gorm.ErrRecordNotFound
}

func (m *invalidPatientUid) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, gorm.ErrRecordNotFound
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}

type MockCal struct{}

func (m *MockCal) CreateEvent(res visit.VisitCalendar) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *MockCal) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *MockCal) UpdateEvent(event *calendar.Event, event_uid string) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *MockCal) DeleteEvent(event_uid string) error {
	return nil
}

type errorCreateEvent struct{}

func (m *errorCreateEvent) CreateEvent(res visit.VisitCalendar) (*calendar.Event, error) {
	return &calendar.Event{}, errors.New("")
}

func (m *errorCreateEvent) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	return &calendar.Event{}, errors.New("")
}

func (m *errorCreateEvent) UpdateEvent(event *calendar.Event, event_uid string) (*calendar.Event, error) {
	return &calendar.Event{}, errors.New("")
}

func (m *errorCreateEvent) DeleteEvent(event_uid string) error {
	return errors.New("")
}

type errorInsertEvent struct{}

func (m *errorInsertEvent) CreateEvent(res visit.VisitCalendar) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *errorInsertEvent) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	return &calendar.Event{}, errors.New("")
}

func (m *errorInsertEvent) UpdateEvent(event *calendar.Event, event_uid string) (*calendar.Event, error) {
	return &calendar.Event{}, errors.New("")
}

func (m *errorInsertEvent) DeleteEvent(event_uid string) error {
	return errors.New("")
}

type errorCancelEvent struct{}

func (m *errorCancelEvent) CreateEvent(res visit.VisitCalendar) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *errorCancelEvent) InsertEvent(event *calendar.Event) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *errorCancelEvent) UpdateEvent(event *calendar.Event, event_uid string) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *errorCancelEvent) DeleteEvent(event_uid string) error {
	return errors.New("")
}

func TestCreate(t *testing.T) {

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
			"doctor_uid": "doctor",
			"date":       "05-05-2022",
			"complaint":  "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  123,
			"patient_uid": "patient",
			"date":        "05-05-2022",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validator doctor_uid", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid": "",
			"date":       "05-05-2022",
			"complaint":  "complaint",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator date", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid": "doctor_uid",
			"date":       "",
			"complaint":  "complaint",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("validator complaint", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid": "doctor_uid",
			"date":       "05-05-2000",
			"complaint":  "",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("invalid date format", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid": "doctor_uid",
			"date":       "05-05-00",
			"complaint":  "comlaint",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("there's another appoinment", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&spesificError{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("leftCapacity", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&leftCapacity{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("doctor_uid", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&invalidDoctorUid{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("patient_uid", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&invalidPatientUid{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error visit list", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&errorVisitList{}, &errorCreateEvent{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 200, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error create event", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorCreateEvent{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 200, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error Insert event", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorInsertEvent{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 200, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error update event id", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
			"complaint":   "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&errorUpdateEventId{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 201, response.Code)
		// log.Info(response.Message)
	})

}

func TestUpdate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/visit")
		context.Param("visit_uid")
		context.SetParamNames("visit_uid")
		context.SetParamValues("visit 123")
		// log.Info(context.ParamNames())

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 202, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": 123,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("date", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"date": "05-05-2002",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&spesificError{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error visit list", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&errorVisitList{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 202, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error create event", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorCreateEvent{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 202, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error update event", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorInsertEvent{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 202, response.Code)
		// log.Info(response.Message)
	})

	t.Run("error delete cancel event event", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
			"status":    "cancelled",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorCancelEvent{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 202, response.Code)
		// log.Info(response.Message)
	})

}

func TestDelete(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/visit")
		context.Param("visit_uid")
		context.SetParamNames("visit_uid")
		context.SetParamValues("visit 123")
		// log.Info(context.ParamNames())

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 202, response.Code)
	})

	t.Run("recordNotFound", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&spesificError{}, &MockCal{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

	t.Run("delete event", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &errorInsertEvent{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 200, response.Code)
		// log.Info(response.Message)
	})

}

func TestGetVisits(t *testing.T) {

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
		context.QueryParams().Add("status", "pending")

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetVisits())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 200, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetVisits())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}
