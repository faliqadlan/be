package visit

import (
	"be/api/google/calendar"
	"be/delivery/controllers/templates"
	logic "be/delivery/logic/visit"
	"be/delivery/middlewares"
	"be/entities"
	"be/repository/visit"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	calGoogle "google.golang.org/api/calendar/v3"
)

type Controller struct {
	r   visit.Visit
	cal calendar.Calendar
	l   logic.Visit
}

func New(r visit.Visit, cal calendar.Calendar, l logic.Visit) *Controller {
	return &Controller{
		r:   r,
		cal: cal,
		l:   l,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {

		var uid string

		if uid = c.QueryParam("patient_uid"); uid == "" {
			uid, _ = middlewares.ExtractTokenUid(c)
		}

		var req logic.Req

		if err := c.Bind(&req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input", nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			log.Warn(err)
			switch {
			case strings.Contains(err.Error(), "Date"):
				err = errors.New("invalid date")
			case strings.Contains(err.Error(), "Doctor_uid"):
				err = errors.New("invalid doctor_uid")
			case strings.Contains(err.Error(), "Complaint"):
				err = errors.New("invalid complaint")
			default:
				err = errors.New("invalid input")
			}
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		if err := cont.l.ValidationRequest(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		entity, err := req.ToVisit()
		if req.Date != "" {
			log.Warn(err)
			if err != nil {
				return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
			}
		}

		res, err := cont.r.CreateVal(req.Doctor_uid, uid, *entity)

		if err != nil {
			log.Warn(err)
			switch {
			case strings.Contains(err.Error(), "doctor_uid"):
				err = errors.New("invalid doctor_uid")
			case strings.Contains(err.Error(), "patient_uid"):
				err = errors.New("invalid patient_uid")
			case err.Error() == errors.New("there's another appoinment in pending").Error():
				err = errors.New("there's another appoinment in pending")
			case err.Error() == errors.New("left capacity can't below zero").Error():
				err = errors.New("left capacity can't below zero")
			default:
				err = errors.New("there's problem in server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		resCal, err := cont.r.GetVisitList(res.Visit_uid)
		if err != nil {
			log.Warn(err)
		}

		var res1 *calGoogle.Event
		// for {
		res1, err = cont.cal.CreateEvent(resCal)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusCreated, templates.Success(nil, "success add visit", res.Complaint))
		}
		res1, err = cont.cal.InsertEvent(res1)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusCreated, templates.Success(nil, "success add visit", res.Complaint))
		}

		_, err = cont.r.Update(res.Visit_uid, entities.Visit{Event_uid: res1.Id})
		if err != nil {
			log.Warn(err)
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add visit and attach to google calendar", res1.HtmlLink))
	}
}

func (cont *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = c.Param("visit_uid")
		var req logic.Req

		if err := c.Bind(&req); err != nil {
			log.Warn(err)
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input", nil))
		}

		switch {
		case req.Date != "":
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "date can't updated, must cancel the appoinment", nil))
		}

		if err := cont.l.ValidationRequest(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, err.Error(), nil))
		}

		// log.Info(uid)
		entity, _ := req.ToVisit()

		res, err := cont.r.Update(uid, *entity)

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("record not found").Error():
				err = errors.New("data is not found")
			default:
				err = errors.New("there's problem in server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		// google calendar
		resCal, err := cont.r.GetVisitList(uid)
		if err != nil {
			log.Warn(err)
		}

		var res1 *calGoogle.Event
		// for {
		res1, err = cont.cal.CreateEvent(resCal)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update visit", res.Complaint))
		}
		res1, err = cont.cal.UpdateEvent(res1, resCal.Event_uid)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update visit", res.Complaint))
		}

		if req.Status == "cancelled" {
			err = cont.cal.DeleteEvent(resCal.Event_uid)
			if err != nil {
				log.Warn(err)
			}
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update visit", res1.HtmlLink))
	}
}

func (cont *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = c.Param("visit_uid")

		resCal, err := cont.r.GetVisitList(uid)
		if err != nil {
			log.Warn(err)
		}

		res, err := cont.r.Delete(uid)

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("record not found").Error():
				err = errors.New("data is not found")
			default:
				err = errors.New("there's problem in server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}
		// google calendar

		err = cont.cal.DeleteEvent(resCal.Event_uid)
		if err != nil {
			log.Warn(err)
			return c.JSON(http.StatusCreated, templates.Success(nil, "success delete visit", res.Complaint))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success delete visit", res.DeletedAt))
	}
}

func (cont *Controller) GetVisits() echo.HandlerFunc {
	return func(c echo.Context) error {
		var kind = c.QueryParam("kind")
		var uid = c.QueryParam("uid")
		var status = c.QueryParam("status")
		var date = c.QueryParam("date")
		var grouped = c.QueryParam("grouped")

		// log.Info(uid, status)

		var res, err = cont.r.GetVisitsVer1(kind, uid, status, date, grouped)

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			default:
				err = errors.New("there's problem in server")
			}
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get list visit", res))
	}
}
