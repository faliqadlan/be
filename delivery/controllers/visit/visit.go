package visit

import (
	"be/api/google/calendar"
	"be/delivery/controllers/templates"
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
}

func New(r visit.Visit, cal calendar.Calendar) *Controller {
	return &Controller{
		r:   r,
		cal: cal,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)
		var req Req

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
		res, err := cont.r.CreateVal(req.Doctor_uid, uid, *req.ToVisit())

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("there's another appoinment in pending").Error():
				err = errors.New("there's another appoinment in pending")
			case errors.New("left capacity can't below zero").Error():
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

		// 	if !strings.Contains(err.Error(), "id") {
		// 		break
		// 	}
		// }
		// log.Info(res1.Id)
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
		var req Req

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid input", nil))
		}

		switch {
		case req.Date != "":
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "date can't updated, must cancel the appoinment", nil))
		}
		// log.Info(uid)
		var res, err = cont.r.Update(uid, *req.ToVisit())

		if err != nil {
			log.Warn(err)
			switch err.Error() {
			case errors.New("record not found").Error():
				err = errors.New("account is not found")
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
				err = errors.New("account is not found")
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
