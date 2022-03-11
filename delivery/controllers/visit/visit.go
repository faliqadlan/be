package visit

import (
	"be/delivery/controllers/templates"
	"be/delivery/middlewares"
	"be/repository/visit"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	r visit.Visit
}

func New(r visit.Visit) *Controller {
	return &Controller{
		r: r,
	}
}

func (cont *Controller) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = middlewares.ExtractTokenUid(c)
		var req Req

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for add visit "+err.Error(), nil))
		}

		var v = validator.New()
		if err := v.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error validator for add visit "+err.(validator.ValidationErrors).Error(), nil))
		}
		var res, err = cont.r.CreateVal(req.Doctor_uid, uid, *req.ToVisit())

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for add visit "+err.Error(), nil))
		}

		return c.JSON(http.StatusCreated, templates.Success(http.StatusCreated, "success add visit", res.Complaint))
	}
}

func (cont *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = c.Param("visit_uid")
		var req Req

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "error binding for update visit "+err.Error(), nil))
		}
		// log.Info(uid)
		var res, err = cont.r.Update(uid, *req.ToVisit())

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for update visit "+err.Error(), nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success update visit", res.Complaint))
	}
}

func (cont *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		var uid = c.Param("visit_uid")

		var res, err = cont.r.Delete(uid)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for delete visit "+err.Error(), nil))
		}

		return c.JSON(http.StatusAccepted, templates.Success(http.StatusAccepted, "success delete visit", res.DeletedAt))
	}
}

func (cont *Controller) GetVisits() echo.HandlerFunc {
	return func(c echo.Context) error {
		var kind = c.QueryParam("kind")
		var uid = c.QueryParam("uid")
		var status = c.QueryParam("status")
		var signStatus = c.QueryParam("sign_status")

		// log.Info(uid, status)

		var res, err = cont.r.GetVisitsVer1(kind, uid, status, signStatus)

		if err != nil {
			// log.Info(err)
			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, "error internal server for get list visit "+err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(http.StatusOK, "success get list visit", res))
	}
}
